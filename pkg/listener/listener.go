/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package listener

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/cloudevents/conformance/pkg/event"
	cfhttp "github.com/cloudevents/conformance/pkg/http"
)

type Listener struct {
	Port    int
	Path    string
	Verbose bool
	Tee     *url.URL

	History int
	ring    *ringBuffer
}

type ringBuffer struct {
	count int
	out   chan event.Event
	mux   sync.Mutex
}

func (r *ringBuffer) Add(ce event.Event) {
	r.mux.Lock()
	defer r.mux.Unlock()

	select {
	// If we can, write to out.
	case r.out <- ce:
		r.count++
		// Done.
	default:
		// If we got blocked, read one from out and write the new event.
		<-r.out
		r.out <- ce
	}
}

func (r *ringBuffer) All() []event.Event {
	all := make([]event.Event, 0, r.count)
	for r.count > 0 {
		r.mux.Lock()

		ce := <-r.out
		all = append(all, ce)
		r.count--

		r.mux.Unlock()
	}
	return all
}

func (l *Listener) Do() error {
	addr := fmt.Sprintf(":%d", l.Port) // TODO: do this listen thing for port 0

	l.ring = &ringBuffer{out: make(chan event.Event, l.History)}

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

	if req.Method == http.MethodGet && req.URL.String() == "/history" {
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

func (l *Listener) ServeHistory(w http.ResponseWriter, _ *http.Request) {
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
