/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package diff

import (
	"errors"
	"fmt"
	cmpdiff "github.com/kylelemons/godebug/diff"
	"io"
	"math"
	"strings"

	"github.com/cloudevents/conformance/pkg/event"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Diff compares two CloudEvents yaml files (or directories) for differences,
// ignoring `Mode` and `TransportExtensions`.
type Diff struct {
	Out io.Writer

	FindBy          []string
	FullDiff        bool
	IgnoreAdditions bool

	FileA string
	FileB string
}

var ignoreOpts = cmpopts.IgnoreFields(event.Event{}, "Mode", "TransportExtensions")

func newTracker(findBy []string, all []event.Event) *tracker {
	return &tracker{
		findBy: findBy,
		all:    all,
		used:   make(map[int]bool, len(all)),
	}
}

type tracker struct {
	findBy []string
	all    []event.Event
	used   map[int]bool
}

type trackerFn func(string, event.Event)

func (t *tracker) findBest(like event.Event, fn trackerFn) (string, error) {
	// Make the label.
	label := ""
	for _, b := range t.findBy {
		if len(label) == 0 {
			label = fmt.Sprintf("%s[%s]", b, like.Get(b))
		} else {
			label = fmt.Sprintf("%s %s[%s]", label, b, like.Get(b))
		}
	}

	best := -1
	bestDiff := math.MaxInt

	// Look for like.
	for i, b := range t.all {
		if t.used[i] {
			continue
		}

		// Test for a match.
		match := true
		for _, key := range t.findBy {
			if like.Get(key) != b.Get(key) {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		// Simple compute the "delta" between line and b, just the length of the diff.
		diff := len(cmp.Diff(like, b, ignoreOpts))
		if diff < bestDiff {
			best = i
			bestDiff = diff
		}
	}
	if best >= 0 {
		fn(label, t.all[best])
		t.used[best] = true
		return label, nil
	}
	return label, io.EOF
}

func (i *Diff) Do() error {
	// Read and parse FileA
	eventsA, err := event.FromYaml(i.FileA, true)
	if err != nil {
		return err
	}

	// Read and parse FileB
	eventsB, err := event.FromYaml(i.FileB, true)
	if err != nil {
		return err
	}

	var diffs []string

	// Tracker helps manage selecting events based on findBy and the current event to compare.
	t := newTracker(i.FindBy, eventsB)

	for _, a := range eventsA {

		if label, err := t.findBest(a, func(label string, b event.Event) {
			if diff := cmp.Diff(a, b, ignoreOpts); diff != "" {
				// Clear out Mode and Transport Extensions.
				a.Mode = ""
				b.Mode = ""
				a.TransportExtensions = nil
				b.TransportExtensions = nil

				ab, err := event.ToYaml(a)
				if err != nil {
					return
				}
				bb, err := event.ToYaml(b)
				if err != nil {
					return
				}

				ignore := true
				sb := &strings.Builder{}
				if len(diffs) >= 1 {
					sb.WriteString("---\n")
				}
				sb.WriteString(fmt.Sprintf("%s diffs (-a, +b):\n", label))

				chunks := cmpdiff.DiffChunks(strings.Split(string(ab), "\n"), strings.Split(string(bb), "\n"))
				changed := map[string]bool{}
				for _, c := range chunks {
					if len(c.Added) != 0 {
						for _, a := range c.Added {
							if _, found := changed[strings.Split(a, ":")[0]]; found || !i.IgnoreAdditions {
								ignore = false
								sb.WriteString(fmt.Sprintf("+ %s\n", a))
							}
						}
					}
					if len(c.Deleted) != 0 {
						ignore = false
						for _, d := range c.Deleted {
							sb.WriteString(fmt.Sprintf("- %s\n", d))
							changed[strings.Split(d, ":")[0]] = true
						}
					}
				}

				if i.FullDiff {
					_, _ = fmt.Fprintf(i.Out, "%s\n", cmpdiff.Diff(string(ab), string(bb)))
				} else if !ignore {
					_, _ = fmt.Fprint(i.Out, sb.String())
					if len(diffs) > 0 {
						diffs = append(diffs, "---")
					}
					diffs = append(diffs, cmpdiff.Diff(string(ab), string(bb)))
				}
			}
		}); err != nil {
			_, _ = fmt.Fprintf(i.Out, "missing: %s\n", label)
			diffs = append(diffs, fmt.Sprintf("missing: %s\n", label))
			continue
		}

	}

	if len(diffs) > 0 {
		return errors.New(strings.Join(diffs, "\n"))
	}

	return nil
}
