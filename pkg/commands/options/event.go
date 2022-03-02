package options

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cloudevents/conformance/pkg/event"
)

// EventOptions
type EventOptions struct {
	Event      event.Event
	Extensions []string // in the form key=value
	Now        bool
}

func wrap80(text string) string {
	return wrap(text, 80)
}

func wrap(text string, width int) string {
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	count := width - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > count {
			wrapped += "\n" + word
			count = width - len(word)
		} else {
			wrapped += " " + word
			count -= 1 + len(word)
		}
	}
	return wrapped
}

const (
	httpSpecMode = `Mode of the outbound event. In the binary content mode, the value of the event data is placed into the HTTP request. In the structured content mode, event metadata attributes and event data are placed into the HTTP request body as JSON. [binary, structured]`
	// Required
	specTextID     = "Identifies the event. Producers MUST ensure that source + id is unique for each distinct event. If a duplicate event is re-sent (e.g. due to a network error) it MAY have the same id. Consumers MAY assume that Events with identical source and id are duplicates."
	specTextSource = "Identifies the context in which an event happened. Often this will include information such as the type of the event source, the organization publishing the event or the process that produced the event. The exact syntax and semantics behind the data encoded in the URI is defined by the event producer."
	specTextType   = "This attribute contains a value describing the type of event related to the originating occurrence. Often this attribute is used for routing, observability, policy enforcement, etc. SHOULD be prefixed with a reverse-DNS name. The prefixed domain dictates the organization which defines the semantics of this event type."
	// Optional
	specTextDataContentType = "Content type of data value. This attribute enables data to carry any type of content, whereby format and encoding might differ from that of the chosen event format."
	specTextDataSchema      = "Identifies the schema that data adheres to. Incompatible changes to the schema SHOULD be reflected by a different URI."
	specTextSubject         = "Describes the subject of the event in the context of the event producer (identified by source). In publish-subscribe scenarios, a subscriber will typically subscribe to events emitted by a source, but the source identifier alone might not be sufficient as a qualifier for any specific event if the source context has internal sub-structure."
	specTextTime            = "Timestamp of when the occurrence happened. MUST adhere to the format specified in RFC 3339."
	specTextExtensions      = "A CloudEvent MAY include any number of additional context attributes with distinct names, known as 'extension attributes'. Extension attributes MUST follow the same naming convention and use the same type system as standard attributes. Extension attributes have no defined meaning in this specification, they allow external systems to attach metadata to an event, much like HTTP custom headers."
	specTextData            = "The event payload. This specification does not place any restriction on the type of this information. It is encoded into a media format which is specified by the datacontenttype attribute (e.g. application/json), and adheres to the dataschema format when those respective attributes are present."
)

func (o *EventOptions) AddFlags(cmd *cobra.Command) {

	// Content
	cmd.Flags().StringVar(&o.Event.Mode, "mode", "", wrap80(httpSpecMode))

	// Required fields.

	// Lock to cloudevents 1.0 for now.
	o.Event.Attributes.SpecVersion = "1.0"

	cmd.Flags().StringVar(&o.Event.Attributes.ID, "id", "", wrap80(specTextID))
	_ = cmd.MarkFlagRequired("id")

	cmd.Flags().StringVar(&o.Event.Attributes.Type, "type", "", wrap80(specTextType))
	_ = cmd.MarkFlagRequired("type")

	cmd.Flags().StringVar(&o.Event.Attributes.Source, "source", "", wrap80(specTextSource))
	_ = cmd.MarkFlagRequired("source")

	// Optional Fields.
	cmd.Flags().StringVar(&o.Event.Attributes.DataContentType, "datacontenttype", "", wrap80(specTextDataContentType))

	cmd.Flags().StringVar(&o.Event.Attributes.DataSchema, "dataschema", "", wrap80(specTextDataSchema))

	cmd.Flags().StringVar(&o.Event.Attributes.Subject, "subject", "", wrap80(specTextSubject))

	cmd.Flags().StringVar(&o.Event.Attributes.Time, "time", "", wrap80(specTextTime))

	cmd.Flags().BoolVar(&o.Now, "timenow", false, "Set time to now.")

	// Extensions Fields.
	cmd.Flags().StringSliceVar(&o.Extensions, "extension", nil, wrap80(specTextExtensions+" Example: key=value."))

	// Data.
	cmd.Flags().StringVar(&o.Event.Data, "data", "", wrap80(specTextData))
}
