/*
 Copyright 2022 The CloudEvents Authors
 SPDX-License-Identifier: Apache-2.0
*/

package model

const (
	Version0dot3 = "v0.3"
	Version1dot0 = "v1.0"
)

type Features struct {
	Language string                       `json:"language" yaml:"language"`
	Versions map[string]VersionedFeatures `json:"versions" yaml:"versions"`
}

type VersionedFeatures struct {
	Core         bool         `json:"core" yaml:"core"`
	EventFormats EventFormats `json:"formats" yaml:"formats"`
	Bindings     Bindings     `json:"bindings" yaml:"bindings"`
}

type EventFormats struct {
	// Avro Event Format
	Avro bool `json:"avro" yaml:"avro"`
	// AMQP Event Format
	AMQP bool `json:"amqp" yaml:"amqp"`
	// JSON Event Format
	JSON bool `json:"json" yaml:"json"`
	// Protobuf Event Format
	Protobuf bool `json:"protobuf" yaml:"protobuf"`
}

type Bindings struct {
	// AMQP Binding
	AMQP ContentModes `json:"amqp" yaml:"amqp"`
	// HTTP Binding
	HTTP ContentModesWithBatch `json:"http" yaml:"http"`
	// Kafka Binding
	Kafka ContentModes `json:"kafka" yaml:"kafka"`
	// MQTT Binding
	MQTT ContentModes `json:"mqtt" yaml:"mqtt"`
	// NATS Binding
	NATS ContentModes `json:"nats" yaml:"nats"`
	// WebSockets Binding
	WebSockets ContentModes `json:"web-sockets" yaml:"web-sockets"`
}

type ContentModes struct {
	// Binary content mode.
	Binary bool `json:"binary" yaml:"binary"`
	// Structured content mode.
	Structured bool `json:"structured" yaml:"structured"`
}

type ContentModesWithBatch struct {
	ContentModes `json:",inline" yaml:",inline"`
	// Batch content mode.
	Batch bool `json:"batch" yaml:"batch"`
}
