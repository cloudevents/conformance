package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func RequestToEvent(req *http.Request) (*event.Event, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = body

	event := &event.Event{
		Data: string(body),
	}

	// CloudEvents attributes.
	event.Attributes.SpecVersion = req.Header.Get("ce-specversion")
	req.Header.Del("ce-specversion")
	event.Attributes.Type = req.Header.Get("ce-type")
	req.Header.Del("ce-type")
	event.Attributes.Time = req.Header.Get("ce-time")
	req.Header.Del("ce-time")
	event.Attributes.ID = req.Header.Get("ce-id")
	req.Header.Del("ce-id")
	event.Attributes.Source = req.Header.Get("ce-source")
	req.Header.Del("ce-source")
	event.Attributes.Subject = req.Header.Get("ce-subject")
	req.Header.Del("ce-subject")
	event.Attributes.SchemaURL = req.Header.Get("ce-schemaurl")
	req.Header.Del("ce-schemaurl")
	event.Attributes.DataContentType = req.Header.Get("Content-Type")
	req.Header.Del("Content-Type")
	event.Attributes.DataContentEncoding = req.Header.Get("ce-datacontentencoding")
	req.Header.Del("ce-datacontentencoding")

	// CloudEvents attribute extensions.
	event.Attributes.Extensions = make(map[string]string)
	for k, _ := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "ce-") {
			event.Attributes.Extensions[k[len("ce-"):]] = req.Header.Get(k)
			req.Header.Del(k)
		}
	}

	// Transport extensions.
	event.TransportExtensions = make(map[string]string)
	for k, _ := range req.Header {
		event.TransportExtensions[k] = req.Header.Get(k)
		req.Header.Del(k)
	}

	return event, nil
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
