package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/cloudevents/conformance/pkg/event"
)

var helloWorldV03 = `
ContextAttributes:
  specversion: 0.3
  type: io.cloudevents.unittest
  id: unit-test-0001
  source: //github.com/cloudevents/conformance/pkg/http/unittest
  datacontenttype: application/json; charset=utf-8
Data: |
  {"msg":"Hello, World!"}
`

func makeUnitTestYaml() (string, func()) {
	file, err := ioutil.TempFile("", "unit_test_*.yaml")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = file.WriteString(helloWorldV03)
	return file.Name(), func() {
		_ = os.Remove(file.Name())
	}
}

func TestHTTP_Binary_v3_yaml(t *testing.T) {
	tmp, done := makeUnitTestYaml()
	defer done()

	server := httptest.NewServer(&dump{t: t})
	events := event.FromYaml(tmp, true)
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
