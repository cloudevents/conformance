package sender

import (
	"fmt"
	"github.com/cloudevents/conformance/pkg/event"
	"github.com/cloudevents/conformance/pkg/http"
	"net/http/httputil"
	"net/url"
)

type Sender struct {
	URL     *url.URL
	Event   event.Event
	Verbose bool
}

func (s *Sender) Do() error {
	req, err := http.EventToRequest(s.URL.String(), s.Event)
	if err != nil {
		return err
	}
	if s.Verbose {
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Printf("Failed to dump request: %+v\n", err)
		} else {
			fmt.Println(string(b))
		}
	}

	if s.URL.Host == "-" {
		// Use "-" as a special hostname to indicate that actual requests should be skipped.
		return nil
	}

	if err := http.Do(req); err != nil {
		return err
	}
	return nil
}
