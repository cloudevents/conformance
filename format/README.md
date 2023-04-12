# Conformance test documentation

The conformance tests for event formats consist of data files in a text representation.
Anyone implementing an event format is strongly encouraged to implement the
conformance tests, which  involves:

- Creating "golden" sample events in code to avoid any ambiguity.
- Implementing a test framework on per-format basis, to load the files
  and process each test within each file.
- Implementing appropriate source control (e.g. a git submodule) to allow the conformance
  tests to be updated on a regular basis. (Implementers are strongly encouraged *not* to
  take a one-off copy of the data files, effectively forking them.)

Any new format proposal should come with a comprehensive set of conformance tests,
implemented in at least one language. Clarifications to a format should also come in the
form of a wording change and corresponding conformance test.

Conformance tests should always be valid in terms of the higher-level format: JSON format
test files are always in the form of valid JSON documents, for example.

## "Golden" sample events

The following sections describe sample events that are referred to by ID in conformance tests.
Each "valid event" or "valid batch" test specifies the expected sample event that the parsed
event is expected to be equal to.

When comparing events for equality, only "CloudEvent-relevant" differences should be reported
as failures. In particular, context attributes with a type of `Timestamp` should only be
compared for equality of "instant in time" represented: where a local offset from UTC is
included, that does not represent a CloudEvent-relevant difference; nor does the precision
of the timestamp. The conformance tests expect timestamp precisions to the millisecond to
be available, but nothing further.

The descriptions below are expected to be implemented in code; they are deliberately *not*
machine-readable, as that would encourage effectively an "extra format" which would be error-prone.

Note that in some event formats, extension attributes may not express type information.
Some SDKs may allow "known" extension attributes to be expressed when parsing, which enables
conformance tests to still check that such extension attributes can be handled appropriately.
Conformance tests which require known sample extension attributes to be provided must indicate
this, and implementations which have no such concept may skip these tests. The list of sample
extension attributes is provided [later in this document](#sample-extension-attributes).

Similarly, some implementations may allow extension attributes to be specified, but without constraints.
Conformance tests which require custom constraints to be applied (i.e. beyond the constraints within
the core specification) must also indicate this requirement, so that implementations without that
feature can skip these tests.

Conformance tests for a format should use these common sample events as far as possible, but may
also define their own sample events to test format-specific aspects such as escaping or particular
data aspects. Format-specific sample events should be prefixed by the format name and a dash,
e.g. "protobuf-messageData".

### minimal

The minimal event is widely used in the conformance tests, and *every* format should have at
least one valid representation. Most other sample events are based on this.

It has the following context attributes, with no optional attributes and no data:

- `specversion`: "1.0"
- `id`: "minimal"
- `type`: "io.cloudevents.test"
- `source`: "https://cloudevents.io"

### minimalWithTime

This event is used for simple validation of aspects of time parsing. It is derived from the minimal event,
with the following changes:

- `id`: "minimalWithTime"
- `time`: 2018-04-05T17:31:00Z

### minimalWithRelativeSource

This event is used for simple validation of URI reference parsing. It is derived from the minimal event,
with the following changes:

- `id`: "minimalWithRelativeSource"
- `source`: "#fragment"

### allCore

This event is used to test all core context attributes. It is derived from the minimal event,
with the following changes:

- `id`: "allCore"
- `datacontenttype`: "text/plain"
- `dataschema`: "https://cloudevents.io/dataschema"
- `subject`: "tests"
- `time`: 2018-04-05T17:31:00Z

Note that despite the presence of the `datacontenttype` attribute, this event has no data.

### simpleTextData

This event is used to test simple data handling. It is derived from the minimal event,
with the following changes:

- `id`: "simpleData"
- `datacontenttype`: "text/plain"
- `data`: "Simple text"

### allExtensionTypes

This event is derived from the minimal event, with the following changes to provide
an extension attribute of each available type:

- `extinteger`: 10
- `extboolean`: true
- `extstring`: "text"
- `extbinary`: bytes `{ 77, 97 }` (or { 0x4d, 0x61 } as hex; base64 is "TWE=")
- `exttimestamp`: 2023-03-31T15:12:00Z
- `exturi`: "https://cloudevents.io"
- `exturiref`: "//authority/path"

## "Golden" sample event batches

### empty

The `empty` batch consists of no events.

### minimal

The `minimal` batch consists of a single event, the `minimal` sample event.

### minimal2

The `minimal2` batch consists of two events, both copies of the `minimal` sample event.
(There is no uniqueness requirement for events in a batch.)

### minimalAndAllCore

The `minimalAndAllCore` batch consists of two events, which are the
`minimal` and then `allCore` sample events.

### minimalAndAllExtensionTypes

The `minimalAndAllExtensionTypes` batch consists of two events, which are the
`minimal` and then `allExtensionTypes` sample events.

## Sample extension attributes

### Unconstrained attributes

- `extinteger`: Integer
- `extboolean`: Boolean
- `extstring`: String
- `extbinary`: Binary
- `exttimestamp`: Timestamp
- `exturi`: URI
- `exturiref`: URI-reference

# Constrained attributes

- `posinteger`: Integer which must be strictly positive
- `shortstring`: String which must have at most 5 characters (conformance tests should only use ASCII here)
- `shortbinary`: Binary which must have at most 5 bytes
- `moderntimestamp`: Timestamp which must be at or after 2000-01-01T00:00:00Z

## TODO: Sample batches
