/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package listener

import (
	"fmt"
	"github.com/cloudevents/conformance/pkg/event"
	cfhttp "github.com/cloudevents/conformance/pkg/http"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Listener struct {
	Port    int
	Path    string
	Verbose bool
	Tee     *url.URL

	History int
	Retain  bool
	ring    *ringBuffer
}

func (l *Listener) Do() error {
	addr := fmt.Sprintf(":%d", l.Port) // TODO: do this listen thing for port 0

	l.ring = newRingBuffer(l.History, l.Retain)

	_, _ = fmt.Fprintf(os.Stderr, "listening on %s\n", addr)
	if err := http.ListenAndServe(addr, l); err != nil {
		return err
	}
	return nil
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if l.Verbose {
		_, _ = fmt.Fprintf(os.Stderr, "incoming request from %s\n", req.URL.String())
	}

	if req.Method == http.MethodGet && req.URL.Path == "/history" {
		l.ServeHistory(w, req)
		return
	}

	ce, err := cfhttp.RequestToEvent(req)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting reqest to event: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if l.History > 0 && ce != nil {
		l.ring.Add(*ce)
	}

	yaml, err := event.ToYaml(*ce)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting event to yaml: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(os.Stdout, string(yaml))
	_, _ = fmt.Fprint(os.Stdout, "---\n")

	if l.Verbose {
		_, _ = fmt.Fprint(os.Stderr, string(yaml))
		_, _ = fmt.Fprint(os.Stderr, "---\n")
	}

	if l.Tee != nil {
		if req, err := cfhttp.EventToRequest(l.Tee.String(), *ce); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error converting event to request: %s\n", err.Error())
		} else if err := cfhttp.Do(req, nil); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error sending event to tee: %s\n", err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (l *Listener) ServeHistory(w http.ResponseWriter, req *http.Request) {
	waitFor := 0
	waitValues, ok := req.URL.Query()["wait"]
	if ok && len(waitValues) == 1 {
		i, err := strconv.Atoi(waitValues[0])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error converting wait value to int: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		waitFor = i
	}
	if waitFor > l.History {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// This will hold the request open until there are at least n events in the ring.
	if waitFor != 0 {
		ticker := time.NewTicker(5 * time.Millisecond)
		wait := true
		for wait {
			select {
			case <-ticker.C:
				if l.ring.Len() >= waitFor {
					wait = false
				}
			case <-req.Context().Done():
				return
			}
		}
	}

	for i, ce := range l.ring.All() {
		yaml, err := event.ToYaml(ce)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error converting event to yaml: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if i > 0 {
			_, _ = w.Write([]byte("---\n"))
		}
		_, _ = w.Write(yaml)
	}
}
