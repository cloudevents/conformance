package listener

import (
	"fmt"
	"github.com/cloudevents/conformance/pkg/event"
	"reflect"
	"testing"
)

func Test_ringBuffer_Add_All(t *testing.T) {
	tests := []struct {
		name  string
		count int
		add   []event.Event
		want  []event.Event
	}{{
		name:  "add 1",
		count: 5,
		add: []event.Event{{
			Mode: "UnitTest1",
		}},
		want: []event.Event{{
			Mode: "UnitTest1",
		}},
	}, {
		name:  "add 5",
		count: 5,
		add: []event.Event{{
			Mode: "UnitTest1",
		}, {
			Mode: "UnitTest2",
		}, {
			Mode: "UnitTest3",
		}, {
			Mode: "UnitTest4",
		}, {
			Mode: "UnitTest5",
		}},
		want: []event.Event{{
			Mode: "UnitTest1",
		}, {
			Mode: "UnitTest2",
		}, {
			Mode: "UnitTest3",
		}, {
			Mode: "UnitTest4",
		}, {
			Mode: "UnitTest5",
		}},
	}, {
		name:  "overflow ring",
		count: 2,
		add: []event.Event{{
			Mode: "UnitTest1",
		}, {
			Mode: "UnitTest2",
		}, {
			Mode: "UnitTest3",
		}, {
			Mode: "UnitTest4",
		}, {
			Mode: "UnitTest5",
		}},
		want: []event.Event{{
			Mode: "UnitTest4",
		}, {
			Mode: "UnitTest5",
		}},
	}}
	for _, tt := range tests {
		for _, retain := range []bool{true, false} {
			t.Run(tt.name+fmt.Sprintf(" retain %v", retain), func(t *testing.T) {
				// Setup
				r := newRingBuffer(tt.count, retain)

				// Add all to ring
				for _, ce := range tt.add {
					r.Add(ce)
				}

				// Test All
				if got := r.All(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("All() = %v, want %v", got, tt.want)
					return
				}

				// Test retain
				if retain {
					if got := r.All(); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("second time, All() = %v, want %v", got, tt.want)
					}
				} else {
					if got := r.All(); !reflect.DeepEqual(got, []event.Event{}) {
						t.Errorf("second time, All() = %v, want %v", got, []event.Event{})
					}
				}
			})
		}
	}
}
