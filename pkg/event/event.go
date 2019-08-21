package event

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Event struct {
	Attributes ContextAttributes          `yaml:"ContextAttributes"`
	Extensions ExtensionContextAttributes `yaml:"ExtensionContextAttributes,omitempty"`
	Data       DataAttributes             `yaml:"DataAttributes"`
}

type ContextAttributes struct {
	SpecVersion         string `yaml:"specversion,omitempty"`
	Type                string `yaml:"type,omitempty"`
	Time                string `yaml:"time,omitempty"`
	ID                  string `yaml:"id,omitempty"`
	Source              string `yaml:"source,omitempty"`
	Subject             string `yaml:"subject,omitempty"`
	SchemaURL           string `yaml:"schemaurl,omitempty"`
	DataContentEncoding string `yaml:"dataecontentncoding,omitempty"`
	DataContentType     string `yaml:"datacontenttype,omitempty"`
}

type ExtensionContextAttributes map[string]string

type DataAttributes struct {
	Data string `yaml:"data,omitempty"`
}

func readFile(pathname string) ([]Event, error) {
	file, err := os.Open(pathname)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	return decode(file)
}

func decode(reader io.Reader) ([]Event, error) {
	decoder := yaml.NewDecoder(reader)
	events := make([]Event, 0)
	var err error
	for {
		out := Event{}
		err = decoder.Decode(&out)
		if err != nil {
			break
		}
		events = append(events, out)
	}
	if err != io.EOF {
		return nil, err
	}
	return events, nil
}

func FromYaml(file string) []Event {
	events, err := readFile(file)
	if err != nil {
		panic(err)
	}
	return events
}

func File(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	root := strings.Split(dir, "conformance")
	return fmt.Sprintf("%s/conformance/yaml/%s", root[0], name)
}
