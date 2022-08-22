/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package renderer

import (
	"html/template"
	"io"

	"github.com/cloudevents/conformance/pkg/features/model"
)

var table = `<!-- START GENERATED TABLE -->
<!-- This table was produced by the feature command from https://github.com/cloudevents/conformance -->
<!-- COMMAND -->
<table>
<thead>
  <tr>
    <th colspan="2">Feature</th>
	{{- with .Language }}{{- range . }}
    <th>{{ . }}</th>
	{{- end }}{{- end }}
  </tr>
</thead>
<tbody>
  <tr>
    <td colspan="2"><strong>v1.0</strong></td>
    <td></td>
  </tr>
  <tr>
    <td colspan="2"><a href="https://github.com/cloudevents/spec/blob/v1.0/spec.md">CloudEvents Core</a></td>
    <td> {{ .true }} </td>
  </tr>
  <tr>
    <td colspan="3"></td>
  </tr>
  {{- range $index, $value = .Versions }}
  <tr>
    <td colspan="2">{{ .Name }}</td>
    <td> {{ .true }} </td>
  </tr>
  {{- end }}{{- end }}

  <tr>
    <td colspan="2">Avro Event Format</td>
    <td> {{ .true }} </td>
  </tr>
  <tr>
    <td colspan="3"></td>
  </tr>
  <tr>
    <td rowspan="3">HTTP</td>
    <td>Binary</td>
    <td> {{ .true }} </td>
  </tr>
  <tr>
    <td>Structured</td>
    <td> {{ .true }} </td>
  </tr>
  <tr>
    <td>Batch</td>
    <td> {{ .false }} </td>
  </tr>
</tbody>
</table>
<!-- END GENERATED TABLE -->
`

var tableold = `<!-- START GENERATED TABLE -->
<!-- This table was produced by the feature command from https://github.com/cloudevents/conformance -->
<!-- COMMAND -->
| Feature | Foo  |
| :------ | :-: |
| **[v1.0](https://github.com/cloudevents/spec/tree/v1.0)** |
| [CloudEvents Core](https://github.com/cloudevents/spec/blob/v1.0/spec.md) | :heavy_check_mark: |
| |
| [Avro Event Format](https://github.com/cloudevents/spec/blob/main/cloudevents/formats/avro-format.md) | :x: |
| [AMQP Event Format - need link]() | :x: |
| [JSON Event Format](https://github.com/cloudevents/spec/blob/main/cloudevents/formats/json-format.md) | :heavy_check_mark: |
| [Protobuf Event Format](https://github.com/cloudevents/spec/blob/main/cloudevents/formats/protobuf-format.md) | :x: |
||
| [AMQP Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/amqp-protocol-binding.md) | :x: |
| [HTTP Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/http-protocol-binding.md) | :heavy_check_mark: |
| [Kafka Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/kafka-protocol-binding.md) :heavy_check_mark: |
| [MQTT Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/mqtt-protocol-binding.md) | :x: |
| [NATS Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/nats-protocol-binding.md) | :x: |
| [WebSockets Binding](https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/websockets-protocol-binding.md) | :x: |
| |
| [HTTP Content Mode Binary](https://github.com/cloudevents/spec/blob/v0.3/http-transport-binding.md) | :heavy_check_mark: |
| [HTTP Content Mode Structured](https://github.com/cloudevents/spec/blob/v0.3/http-transport-binding.md) | :heavy_check_mark: |
| [HTTP Content Mode Batch](https://github.com/cloudevents/spec/blob/v0.3/http-transport-binding.md) | :heavy_check_mark: |
| [Kafka Content Mode Structured]() | :heavy_check_mark: |
| [Kafka Content Mode Binary]() | :heavy_check_mark: |
| [Kafka Content Mode Structured]() | :heavy_check_mark: |
| [MQTT Content Mode Binary](https://github.com/cloudevents/spec/blob/v0.3/mqtt-transport-binding.md) | :heavy_check_mark: |
| [MQTT Content Mode Structured](https://github.com/cloudevents/spec/blob/v0.3/mqtt-transport-binding.md) | :heavy_check_mark: |
| |
<!-- END GENERATED TABLE -->
`

type Renderer struct {
}

func (r *Renderer) Render(input []model.Features, out io.Writer) error {
	tbl, err := template.New("table").Parse(table)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
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

	// TODO: likely have to process input into a flat struct.
	return tbl.Execute(out, data)
}
