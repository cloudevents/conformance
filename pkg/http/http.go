package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudevents/conformance/pkg/event"
)

func addHeader(req *http.Request, key, value string) {
	value = strings.TrimSpace(value)
	if value != "" {
		req.Header.Add(key, value)
	}
}

func EventToRequest(url string, event event.Event) (*http.Request, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(event.Data)))
	if err != nil {
		return nil, err
	}

	// CloudEvents attributes.
	addHeader(req, "ce-specversion", event.Attributes.SpecVersion)
	addHeader(req, "ce-type", event.Attributes.Type)
	addHeader(req, "ce-time", event.Attributes.Time)
	addHeader(req, "ce-id", event.Attributes.ID)
	addHeader(req, "ce-source", event.Attributes.Source)
	addHeader(req, "ce-subject", event.Attributes.Subject)
	addHeader(req, "ce-schemaurl", event.Attributes.SchemaURL)
	addHeader(req, "Content-Type", event.Attributes.DataContentType)
	addHeader(req, "ce-datacontentencoding", event.Attributes.DataContentEncoding)

	// CloudEvents attribute extensions.
	for k, v := range event.Attributes.Extensions {
		addHeader(req, "ce-"+k, v)
	}

	// Transport extensions.
	for k, v := range event.TransportExtensions {
		addHeader(req, k, v)
	}

	return req, nil
}

func Do(req *http.Request) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("expected 200 level response, got %s", resp.Status)
	}

	// TODO might want something from resp.
	return nil
}
