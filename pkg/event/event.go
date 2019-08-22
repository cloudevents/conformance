package event

type Event struct {
	Attributes          ContextAttributes `yaml:"ContextAttributes"`
	TransportExtensions Extensions        `yaml:"TransportExtensions,omitempty"`
	Data                string            `yaml:"Data"`
}

type ContextAttributes struct {
	SpecVersion         string     `yaml:"specversion,omitempty"`
	Type                string     `yaml:"type,omitempty"`
	Time                string     `yaml:"time,omitempty"`
	ID                  string     `yaml:"id,omitempty"`
	Source              string     `yaml:"source,omitempty"`
	Subject             string     `yaml:"subject,omitempty"`
	SchemaURL           string     `yaml:"schemaurl,omitempty"`
	DataContentEncoding string     `yaml:"dataecontentncoding,omitempty"`
	DataContentType     string     `yaml:"datacontenttype,omitempty"`
	Extensions          Extensions `yaml:"Extensions,omitempty"`
}

type Extensions map[string]string
