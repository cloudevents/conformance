package http

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/cloudevents/conformance/pkg/event"
)

func TestHTTP_Binary_v3_yaml(t *testing.T) {
	server := httptest.NewServer(&dump{t: t})
	events := event.FromYaml(event.File("/"), true)
	testServer := server.URL

	for _, e := range events {
		req, err := EventToRequest(testServer, e)
		if err != nil {
			t.Fatal(err)
		}

		if err := Do(req); err != nil {
			t.Fatal(err)
		}
	}
}

type dump struct {
	t *testing.T
}

func (d *dump) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if reqBytes, err := httputil.DumpRequest(r, true); err == nil {
		d.t.Logf("Received a message: %+v", string(reqBytes))
		_, _ = w.Write(reqBytes)
	} else {
		d.t.Logf("Error: %+v :: %+v", err, r)
	}
}
