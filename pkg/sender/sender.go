package sender

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/cloudevents/conformance/pkg/event"
	"github.com/cloudevents/conformance/pkg/http"
)

type Sender struct {
	URL     *url.URL
	Event   event.Event
	YAML    bool
	Verbose bool

	// PreHook allows for mutation of the outbound event before translation to
	// a to HTTP request.
	PreHook event.MutationFn
	// PostHook allows for recording of the outbound HTTP request and resulting
	// response and/or error.
	PostHook http.ResultsFn
}

func (s *Sender) Do() error {
	var err error
	e := s.Event
	if s.PreHook != nil {
		e, err = s.PreHook(e)
		if err != nil {
			return err
		}
	}

	req, err := http.EventToRequest(s.URL.String(), e)
	if err != nil {
		return err
	}
	if s.Verbose {
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to dump request: %+v\n", err)
		} else {
			_, _ = fmt.Fprint(os.Stderr, string(b))
		}
	}

	yaml, err := event.ToYaml(e)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error converting event to yaml: %s\n", err.Error())
	} else if s.YAML {
		_, _ = fmt.Fprint(os.Stdout, string(yaml))
		_, _ = fmt.Fprint(os.Stdout, "---\n")
	}

	if s.URL.Host == "-" {
		// Use "-" as a special hostname to indicate that actual requests should be skipped.
		return nil
	}

	if err := http.Do(req, s.PostHook); err != nil {
		return err
	}
	return nil
}
