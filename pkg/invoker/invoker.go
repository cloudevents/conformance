package invoker

import (
	"errors"
	"strings"

	"github.com/cloudevents/conformance/pkg/event"
	"github.com/cloudevents/conformance/pkg/http"
)

type Invoker struct {
	URL       string
	Files     []string
	Recursive bool
}

func (i *Invoker) Do() error {
	events := event.FromYaml(strings.Join(i.Files, ","), i.Recursive)

	var errs = make([]string, 0)

	for _, e := range events {
		req, err := http.EventToRequest(i.URL, e)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		if err := http.Do(req); err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
