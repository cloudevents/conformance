package invoker

import (
	"errors"
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/cloudevents/conformance/pkg/event"
	"github.com/cloudevents/conformance/pkg/http"
)

type Invoker struct {
	URL       *url.URL
	Files     []string
	Delay     *time.Duration
	Recursive bool
	Verbose   bool
}

func (i *Invoker) Do() error {
	events := event.FromYaml(strings.Join(i.Files, ","), i.Recursive)

	var errs = make([]string, 0)

	for _, e := range events {
		req, err := http.EventToRequest(i.URL.String(), e)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if i.Verbose {
			b, err := httputil.DumpRequestOut(req, true)
			if err != nil {
				fmt.Printf("Failed to dump request: %+v\n", err)
			} else {
				fmt.Println(string(b))
			}
		}

		if i.URL.Host == "-" {
			// Use "-" as a special hostname to indicate that actual requests should be skipped.
			continue
		}

		if err := http.Do(req); err != nil {
			errs = append(errs, err.Error())
			continue
		}

		if i.Delay != nil {
			time.Sleep(*i.Delay)
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
