/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package diff

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudevents/conformance/pkg/event"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	cmpdiff "github.com/kylelemons/godebug/diff"
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

func (i *Diff) find(like event.Event, all []event.Event) (string, *event.Event) {
	label := ""
	for _, b := range i.FindBy {
		if len(label) == 0 {
			label = fmt.Sprintf("%s[%s]", b, like.Get(b))
		} else {
			label = fmt.Sprintf("%s %s[%s]", label, b, like.Get(b))
		}
	}

	// Look for the matching ID, Source and Subject.
	for _, a := range all {
		match := true
		for _, key := range i.FindBy {
			if a.Get(key) != like.Get(key) {
				match = false
				break
			}
		}
		if match {
			return label, &a
		}
	}
	return label, nil
}

func (i *Diff) Do() error {
	ignoreOpts := cmpopts.IgnoreFields(event.Event{}, "Mode", "TransportExtensions")

	eventsA, err := event.FromYaml(i.FileA, true)
	if err != nil {
		return err
	}

	eventsB, err := event.FromYaml(i.FileB, true)
	if err != nil {
		return err
	}

	var diffs []string

	for _, a := range eventsA {
		label, b := i.find(a, eventsB)

		if b == nil {
			_, _ = fmt.Fprintf(i.Out, "missing: %s\n", label)
			diffs = append(diffs, fmt.Sprintf("missing: %s\n", label))
			continue
		}

		if diff := cmp.Diff(a, b, ignoreOpts); diff != "" {
			// Clear out Mode and Transport Extensions.
			a.Mode = ""
			b.Mode = ""
			a.TransportExtensions = nil
			b.TransportExtensions = nil

			ab, err := event.ToYaml(a)
			if err != nil {
				return err
			}
			bb, err := event.ToYaml(*b)
			if err != nil {
				return err
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
	}

	if len(diffs) > 0 {
		return errors.New(strings.Join(diffs, "\n"))
	}

	return nil
}
