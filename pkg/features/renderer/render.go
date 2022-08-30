/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package renderer

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/cloudevents/conformance/pkg/features/model"
)

var table = `
{{- $true := glyph true }}
{{- $false := glyph false }}
<table>
<thead>
  <tr>
    <th colspan="2">Feature</th>
	{{- range .Languages }}
    <th>{{ . }}</th>
	{{- end }}
  </tr>
</thead>
<tbody>
{{- range $vi, $version := .Versions }}
  {{- if gt $vi 0 }}
  <tr>
    <td colspan="{{ fullSpan }}"></td>
  </tr>
  {{- end }}
  <tr>
    <td colspan="2"><strong>{{ $version.Name }}</strong></td>
    {{- range $version.Languages }}
    <th></th>
	{{- end }}
  </tr>
  <tr>
    <td colspan="2"><a href="https://github.com/cloudevents/spec/blob/{{ .Name }}/spec.md">CloudEvents Core</a></td>
    {{- range $version.Core }}
    <td>{{ glyph . }}</td>
	{{- end }}
  </tr>

  <tr>
    <td colspan="{{ fullSpan }}"></td>
  </tr>

  {{- range $version.Formats }}
  <tr>
    <td colspan="2">{{ .Name }}</td>
    {{- range .Support }}
    <td>{{ glyph . }}</td>
	{{- end }}
  </tr>
  {{- end }}

  <tr>
    <td colspan="{{ fullSpan }}"></td>
  </tr>

  {{- range $bi, $binding := $version.Bindings }}{{- range $index, $format := $binding.Formats }}
  <tr>
    {{- if eq $index 0 }}
    <td rowspan="{{ len $binding.Formats }}">{{ $binding.Name }}</td>
    {{- end }}
    <td>{{ $format.Name }}</td>
    {{- range $format.Support }}
    <td>{{ glyph . }}</td>
	{{- end }}
  </tr>
  {{- end }}{{- end }}

{{- end }}
</tbody>
</table>
`

type Renderer struct {
	Languages []string
	Versions  []RenderVersion
}

func (r *Renderer) Version(name string) *RenderVersion {
	for i, v := range r.Versions {
		if v.Name == name {
			return &r.Versions[i]
		}
	}
	// Does not exist.
	r.Versions = append(r.Versions, RenderVersion{
		Name:      name,
		Languages: r.Languages,
		Core:      make([]bool, len(r.Languages)),
	})
	return r.Version(name)
}

func (r *RenderVersion) Format(name string) *RenderFormat {
	for i, f := range r.Formats {
		if f.Name == name {
			return &r.Formats[i]
		}
	}
	// Does not exist.
	r.Formats = append(r.Formats, RenderFormat{
		Name:    name,
		Support: make([]bool, len(r.Languages)),
	})
	return r.Format(name)
}

func (r *RenderVersion) Binding(name string) *RenderBindings {
	for i, f := range r.Bindings {
		if f.Name == name {
			return &r.Bindings[i]
		}
	}
	// Does not exist.
	r.Bindings = append(r.Bindings, RenderBindings{
		Name:      name,
		Languages: r.Languages,
	})
	return r.Binding(name)
}

type RenderVersion struct {
	Name      string
	Languages []string
	Core      []bool
	Formats   []RenderFormat
	Bindings  []RenderBindings
}

type RenderFormat struct {
	Name    string
	Support []bool
}

type RenderBindings struct {
	Name      string
	Languages []string
	Formats   []RenderBindingFormat
}

func (r *RenderBindings) Format(name string) *RenderBindingFormat {
	for i, b := range r.Formats {
		if b.Name == name {
			return &r.Formats[i]
		}
	}
	// Does not exist.
	r.Formats = append(r.Formats, RenderBindingFormat{
		Name:    name,
		Support: make([]bool, len(r.Languages)),
	})
	return r.Format(name)
}

type RenderBindingFormat struct {
	Name    string
	Support []bool
}

func (r *Renderer) Render(input []model.Features, out io.Writer) error {
	if len(input) == 0 {
		return errors.New("need at least one features config")
	}
	tbl, err := template.New("table").Funcs(template.FuncMap{
		"glyph": func(t bool) string {
			if t {
				return ":heavy_check_mark:"
			}
			return ":x:"
		},
		"fullSpan": func() string {
			return strconv.Itoa(2 + len(input))
		},
	}).Parse(table)
	if err != nil {
		return err
	}

	data := Renderer{
		Languages: make([]string, len(input)),
	}

	for i, f := range input {
		data.Languages[i] = f.Language
		for name, vf := range f.Versions {
			v := data.Version(name)
			v.Core[i] = vf.Core

			// The keys below become the display text on in the rendered table.

			v.Format("Avro Event Format").Support[i] = vf.EventFormats.Avro
			v.Format("AMQP Event Format").Support[i] = vf.EventFormats.AMQP
			v.Format("JSON Event Format").Support[i] = vf.EventFormats.JSON
			v.Format("Protobuf Event Format").Support[i] = vf.EventFormats.Protobuf

			v.Binding("AMQP").Format("Binary").Support[i] = vf.Bindings.AMQP.Binary
			v.Binding("AMQP").Format("Structured").Support[i] = vf.Bindings.AMQP.Structured

			v.Binding("HTTP").Format("Binary").Support[i] = vf.Bindings.HTTP.Binary
			v.Binding("HTTP").Format("Structured").Support[i] = vf.Bindings.HTTP.Structured
			v.Binding("HTTP").Format("Batch").Support[i] = vf.Bindings.HTTP.Batch

			v.Binding("Kafka").Format("Binary").Support[i] = vf.Bindings.Kafka.Binary
			v.Binding("Kafka").Format("Structured").Support[i] = vf.Bindings.Kafka.Structured

			v.Binding("MQTT").Format("Binary").Support[i] = vf.Bindings.MQTT.Binary
			v.Binding("MQTT").Format("Structured").Support[i] = vf.Bindings.MQTT.Structured

			v.Binding("NATS").Format("Binary").Support[i] = vf.Bindings.NATS.Binary
			v.Binding("NATS").Format("Structured").Support[i] = vf.Bindings.NATS.Structured

			v.Binding("Web Sockets").Format("Binary").Support[i] = vf.Bindings.WebSockets.Binary
			v.Binding("Web Sockets").Format("Structured").Support[i] = vf.Bindings.WebSockets.Structured
		}
	}

	// TODO: URL support.

	data1 := map[string]interface{}{
		"true": ":heavy_check_mark:", "false": ":x:",
		"Language": []string{"Foo", "Bar", "Baz"},
		"Versions": []map[string]interface{}{
			{"v1.0": map[string]interface{}{
				"Formats": []map[string]interface{}{
					{"Name": "Avro Event Format", "URL": "https://github.com/cloudevents/spec/blob/main/cloudevents/formats/avro-format.md"},
					{"Name": "AMQP Event Format", "URL": "need link"},
					{"Name": "JSON Event Format", "URL": "https://github.com/cloudevents/spec/blob/main/cloudevents/formats/json-format.md"},
					{"Name": "Protobuf Event Format", "URL": "https://github.com/cloudevents/spec/blob/main/cloudevents/formats/protobuf-format.md"},
				},
			}},
		},
	}
	_ = data1

	_, _ = fmt.Fprintf(out, "<!-- START GENERATED TABLE -->\n")
	_, _ = fmt.Fprintf(out, "<!-- This table was produced by the feature command from https://github.com/cloudevents/conformance -->\n")
	_, _ = fmt.Fprintf(out, "<!-- via: %s -->\n", strings.Join(os.Args, " "))

	// TODO: likely have to process input into a flat struct.
	if err := tbl.Execute(out, data); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(out, "<!-- END GENERATED TABLE -->\n")
	return nil
}
