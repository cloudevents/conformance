package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudevents/conformance/pkg/event"
)

type ResultsFn func(*http.Request, *http.Response, error)

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

	out := &event.Event{
		Mode: event.BinaryMode,
		Data: string(body),
	}

	// CloudEvents attributes.
	out.Attributes.SpecVersion = req.Header.Get("ce-specversion")
	req.Header.Del("ce-specversion")
	out.Attributes.Type = req.Header.Get("ce-type")
	req.Header.Del("ce-type")
	out.Attributes.Time = req.Header.Get("ce-time")
	req.Header.Del("ce-time")
	out.Attributes.ID = req.Header.Get("ce-id")
	req.Header.Del("ce-id")
	out.Attributes.Source = req.Header.Get("ce-source")
	req.Header.Del("ce-source")
	out.Attributes.Subject = req.Header.Get("ce-subject")
	req.Header.Del("ce-subject")
	out.Attributes.SchemaURL = req.Header.Get("ce-schemaurl")
	req.Header.Del("ce-schemaurl")
	out.Attributes.DataContentType = req.Header.Get("Content-Type")
	req.Header.Del("Content-Type")
	out.Attributes.DataContentEncoding = req.Header.Get("ce-datacontentencoding")
	req.Header.Del("ce-datacontentencoding")

	// CloudEvents attribute extensions.
	out.Attributes.Extensions = make(map[string]string)
	for k := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "ce-") {
			out.Attributes.Extensions[k[len("ce-"):]] = req.Header.Get(k)
			req.Header.Del(k)
		}
	}

	// Transport extensions.
	out.TransportExtensions = make(map[string]string)
	for k := range req.Header {
		out.TransportExtensions[k] = req.Header.Get(k)
		req.Header.Del(k)
	}

	return out, nil
}

func Do(req *http.Request, hook ResultsFn) error {
	resp, err := http.DefaultClient.Do(req)

	if hook != nil {
		// Non-blocking.
		go hook(req, resp, err)
	}

	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("expected 200 level response, got %s", resp.Status)
	}

	// TODO might want something from resp.
	return nil
}
