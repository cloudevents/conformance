package diff

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestDiff_Do(t *testing.T) {
	tests := []struct {
		name            string
		a               string
		b               string
		findBy          []string
		ignoreAdditions bool
		wantErr         bool
		wantDiff        string
	}{{
		name:    "sample1: two files",
		a:       "./testdata/sample1_a.yaml",
		b:       "./testdata/sample1_b.yaml",
		findBy:  []string{"id", "source"},
		wantErr: false,
	}, {
		name:    "sample2: one folder, one file",
		a:       "./testdata/sample2_a/",
		b:       "./testdata/sample2_b.yaml",
		findBy:  []string{"id", "source"},
		wantErr: false,
	}, {
		name:    "sample3: has diff",
		a:       "./testdata/sample3_a.yaml",
		b:       "./testdata/sample3_b.yaml",
		findBy:  []string{"id", "source"},
		wantErr: true,
		wantDiff: `id[4321-4321-4321-a] source[/mycontext/subcontext] diffs (-a, +b):
-   type: com.example.someevent
+   type: com.example.some.other.event
`,
	}, {
		name:            "sample4: just look at type",
		a:               "./testdata/sample4_a.yaml",
		b:               "./testdata/sample4_b.yaml",
		findBy:          []string{"type"},
		ignoreAdditions: true,
		wantErr:         false,
	}, {
		name:            "sample5: just type, from a dir",
		a:               "./testdata/sample5/want",
		b:               "./testdata/sample5/got",
		findBy:          []string{"type"},
		ignoreAdditions: true,
		wantErr:         false,
	}, {
		name:            "sample6: from a dir with diffs",
		a:               "./testdata/sample6/want",
		b:               "./testdata/sample6/got",
		findBy:          []string{"type"},
		ignoreAdditions: true,
		wantErr:         true,
		wantDiff: `missing: type[com.example.someevent.b]
---
type[com.example.someevent.c] diffs (-a, +b):
-   source: /mycontext/subcontext
+   source: /mycontext/subcontext2
`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			i := &Diff{
				Out:             out,
				FindBy:          tt.findBy,
				IgnoreAdditions: tt.ignoreAdditions,
				FileA:           tt.a,
				FileB:           tt.b,
			}
			if err := i.Do(); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := out.String()
			t.Logf(got)
			if want := tt.wantDiff; want != got {
				t.Fatalf("Found diffs (-want, +got): %s", cmp.Diff(want, got))
			}
		})
	}
}
