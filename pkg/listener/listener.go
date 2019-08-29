package listener

import (
	"fmt"
	"github.com/cloudevents/conformance/pkg/event"
	cfhttp "github.com/cloudevents/conformance/pkg/http"
	"net/http"
	"os"
)

type Listener struct {
	Port    int
	Path    string
	Verbose bool
}

func (l *Listener) Do() error {
	addr := fmt.Sprintf(":%d", l.Port) // TODO: do this listen thing for port 0

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

	ce, err := cfhttp.RequestToEvent(req)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting reqest to event: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
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

	w.WriteHeader(http.StatusOK)
}
