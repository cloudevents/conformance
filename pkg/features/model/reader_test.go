/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package model

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFromYaml(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		recursive bool
		want      []Features
		wantErr   bool
	}{{
		name:      "only http",
		files:     []string{"./testdata/only-http.yaml"},
		recursive: false,
		want: []Features{{
			Language: "Only HTTP",
			Versions: map[string]VersionedFeatures{"v1.0": {
				Core: true,
				EventFormats: EventFormats{
					JSON: true,
				},
				Bindings: Bindings{
					HTTP: ContentModesWithBatch{
						ContentModes: ContentModes{
							Binary:     true,
							Structured: true,
						},
						Batch: true,
					},
				},
			},
			},
		}},
		wantErr: false,
	}, {
		name:      "only http + only proto",
		files:     []string{"./testdata/only-http.yaml", "./testdata/only-proto.yaml"},
		recursive: false,
		want: []Features{{
			Language: "Only HTTP",
			Versions: map[string]VersionedFeatures{"v1.0": {
				Core: true,
				EventFormats: EventFormats{
					JSON: true,
				},
				Bindings: Bindings{
					HTTP: ContentModesWithBatch{
						ContentModes: ContentModes{
							Binary:     true,
							Structured: true,
						},
						Batch: true,
					},
				},
			}},
		}, {
			Language: "Only Proto",
			Versions: map[string]VersionedFeatures{"v1.0": {
				EventFormats: EventFormats{
					Protobuf: true,
				},
			}},
		}},
		wantErr: false,
	}, {
		name:      "all supported v1.0",
		files:     []string{"./testdata/all-supported-v1.0.yaml"},
		recursive: false,
		want: []Features{{
			Language: "All Supported v1.0",
			Versions: map[string]VersionedFeatures{"v1.0": {
				Core: true,
				EventFormats: EventFormats{
					Avro:     true,
					AMQP:     true,
					JSON:     true,
					Protobuf: true,
				},
				Bindings: Bindings{
					AMQP: ContentModes{
						Binary:     true,
						Structured: true,
					},
					HTTP: ContentModesWithBatch{
						ContentModes: ContentModes{
							Binary:     true,
							Structured: true,
						},
						Batch: true,
					},
					Kafka: ContentModes{
						Binary:     true,
						Structured: true,
					},
					MQTT: ContentModes{
						Binary:     true,
						Structured: true,
					},
					NATS: ContentModes{
						Binary:     true,
						Structured: true,
					},
					WebSockets: ContentModes{
						Binary:     true,
						Structured: true,
					},
				},
			}},
		}},
		wantErr: false,
	}, {
		name:      "all supported v0.3 and v1.0",
		files:     []string{"./testdata/all-supported-v0.3-v1.0.yaml"},
		recursive: false,
		want: []Features{{
			Language: "All Supported v0.3, v1.0",
			Versions: map[string]VersionedFeatures{
				"v0.3": {
					Core: true,
					EventFormats: EventFormats{
						Avro:     true,
						AMQP:     true,
						JSON:     true,
						Protobuf: true,
					},
					Bindings: Bindings{
						AMQP: ContentModes{
							Binary:     true,
							Structured: true,
						},
						HTTP: ContentModesWithBatch{
							ContentModes: ContentModes{
								Binary:     true,
								Structured: true,
							},
							Batch: true,
						},
						Kafka: ContentModes{
							Binary:     true,
							Structured: true,
						},
						MQTT: ContentModes{
							Binary:     true,
							Structured: true,
						},
						NATS: ContentModes{
							Binary:     true,
							Structured: true,
						},
						WebSockets: ContentModes{
							Binary:     true,
							Structured: true,
						},
					},
				},
				"v1.0": {
					Core: true,
					EventFormats: EventFormats{
						Avro:     true,
						AMQP:     true,
						JSON:     true,
						Protobuf: true,
					},
					Bindings: Bindings{
						AMQP: ContentModes{
							Binary:     true,
							Structured: true,
						},
						HTTP: ContentModesWithBatch{
							ContentModes: ContentModes{
								Binary:     true,
								Structured: true,
							},
							Batch: true,
						},
						Kafka: ContentModes{
							Binary:     true,
							Structured: true,
						},
						MQTT: ContentModes{
							Binary:     true,
							Structured: true,
						},
						NATS: ContentModes{
							Binary:     true,
							Structured: true,
						},
						WebSockets: ContentModes{
							Binary:     true,
							Structured: true,
						},
					},
				},
			},
		}},
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromYaml(tt.recursive, tt.files...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Found diffs (-want, +got): %s", cmp.Diff(tt.want, got))
			}
		})
	}
}
