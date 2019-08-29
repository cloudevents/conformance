package event

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func FromYaml(files string, recursive bool) []Event {
	pathNames := strings.Split(files, ",")
	events := make([]Event, 0)
	for _, pathName := range pathNames {
		var event []Event
		var err error
		if pathName == "-" {
			event, err = decode(os.Stdin)
		} else {
			event, err = read(pathName, recursive)
		}
		if err != nil {
			panic(err) // TODO: could not panic later.
		}
		events = append(events, event...)
	}
	return events
}

func read(pathName string, recursive bool) ([]Event, error) {
	info, err := os.Stat(pathName)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return readDir(pathName, recursive)
	}
	return readFile(pathName)
}

func readFile(pathName string) ([]Event, error) {
	file, err := os.Open(pathName)
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

func readDir(pathName string, recursive bool) ([]Event, error) {
	list, err := ioutil.ReadDir(pathName)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)
	for _, f := range list {
		name := path.Join(pathName, f.Name())
		var evs []Event

		switch {
		case f.IsDir() && recursive:
			evs, err = readDir(name, recursive)
		case !f.IsDir():
			evs, err = readFile(name)
		}

		if err != nil {
			return nil, err
		}
		events = append(events, evs...)
	}
	return events, nil
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
