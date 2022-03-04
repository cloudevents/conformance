/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package listener

import (
	"sync"

	"github.com/cloudevents/conformance/pkg/event"
)

func newRingBuffer(size int, retain bool) *ringBuffer {
	return &ringBuffer{
		out:    make(chan event.Event, size),
		retain: retain,
	}
}

type ringBuffer struct {
	count  int
	retain bool
	out    chan event.Event
	mux    sync.Mutex
}

func (r *ringBuffer) Add(ce event.Event) {
	r.mux.Lock()
	defer r.mux.Unlock()

	select {
	// If we can, write to out.
	case r.out <- ce:
		r.count++
		// Done.
	default:
		// If we got blocked, read one from out and write the new event.
		<-r.out
		r.out <- ce
	}
}

func (r *ringBuffer) Len() int {
	return r.count
}

func (r *ringBuffer) All() []event.Event {
	r.mux.Lock()
	defer r.mux.Unlock()

	// All captures the expected length of the ring.
	all := make([]event.Event, 0, r.count)
	for i := 0; i < cap(all); i++ {
		ce := <-r.out
		all = append(all, ce)
		if r.retain {
			// Add ce back to the channel.
			r.out <- ce
		} else {
			r.count--
		}
	}
	return all
}
