/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/cloudevents/conformance/pkg/diff"
	"github.com/cloudevents/conformance/pkg/event"
	"github.com/cloudevents/conformance/pkg/listener"
	"github.com/cloudevents/conformance/pkg/sender"
)

// TestDiff produces several events and then compares what was collected with
// what is expected.
func TestDiff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	port := 52318 // could be 0, but this causes a race in testing.

	// Stand up a listener.
	l := &listener.Listener{
		Port:    port,
		Path:    "/",
		Verbose: false,
		History: 10,
		Retain:  false,
	}
	go func() {
		if err := l.Do(ctx); !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("listerer failed to close cleanly, %v", err)
		}
	}()

	// Let server start.
	time.Sleep(time.Millisecond * 10)

	expectFile, err := os.CreateTemp(os.TempDir(), "ce*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("expect file:", expectFile.Name())
	defer os.Remove(expectFile.Name())

	collectFile, err := os.CreateTemp(os.TempDir(), "ce*")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("collect file:", expectFile.Name())
	defer os.Remove(collectFile.Name())

	// Send n events.
	for n := 0; n < 5; n++ {
		s := &sender.Sender{
			URL: &url.URL{
				Scheme: "http", Host: fmt.Sprintf("localhost:%d", port),
			},
			Event: event.Event{
				Mode: "binary",
				Attributes: event.ContextAttributes{
					SpecVersion: "1.0",
					Type:        fmt.Sprintf("unit.test.%d", n),
					ID:          fmt.Sprintf("abc.123.%d", n),
					Source:      "http://unit.test",
				},
			},
			YAML: true,
		}
		if err := s.Do(); err != nil {
			t.Errorf("sender failed to Do(), %s", err)
		}

		expect, err := event.ToYaml(event.Event{
			// We are only going to diff on source and type for the test.
			Attributes: event.ContextAttributes{
				Type:   fmt.Sprintf("unit.test.%d", n),
				Source: "http://unit.test",
			},
		})
		if err != nil {
			t.Errorf("failed to convert expected event to yaml, %v", err)
		}
		if n > 0 {
			_, _ = expectFile.WriteString("---\n")
		}
		// Store each sent event structure as yaml.
		_, _ = expectFile.Write(expect)
	}
	_ = expectFile.Close()

	// Collect and store history from listener.
	target := &url.URL{
		Scheme: "http", Host: fmt.Sprintf("localhost:%d", l.Port), Path: "/history",
	}
	history, err := curl(target.String())
	if err != nil {
		t.Errorf("failed to curl listener history, %v", err)
	}
	_, _ = collectFile.Write(history)
	_ = collectFile.Close()

	fmt.Println("got history: ", string(history))

	// Compute diff between what was sent and what was collected.
	d := &diff.Diff{
		Out:             os.Stdout,
		FindBy:          []string{"type"},
		FullDiff:        false,
		IgnoreAdditions: true,
		FileA:           expectFile.Name(),
		FileB:           collectFile.Name(),
	}

	if err := d.Do(); err != nil {
		t.Errorf("failed do diff, %v", err)
	}

	cancel()
}

func curl(target string) ([]byte, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
