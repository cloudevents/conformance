# XML format conformance tests

Each test file is a valid XML document. The root element is always `<testFile>` with
the following attribute:

- `testType`: the default test type for all tests in this file;
  individual tests may override this. Values:
  - `ValidSingleEvent`
  - `ValidBatch`
  - `InvalidSingleEvent`
  - `InvalidBatch`

The root element contains child `<test>` elements, each of which has the following
attributes:

- `id`: the ID of the test
- `sampleId`: for a valid event/batch test, this provides the "golden"
  sample ID to compare the result with. This may be absent, in order
  to test valid events for which there is no golden sample.
- `testType`: used to override the default test type in the file
- `description`: an optional description of the test (purely for clarity)
- `sampleExtensionAttributes`: `true` if this test requires pre-defined sample extension attributes;
  implementations which do not support providing extension attribute definitions during
  parsing may skip these tests.
- `extensionConstraints`: `true` if this test requires constraints to be applied;
  implementations which do not apply constraints may skip these tests.

Each test element has exactly one child element, representing the batch or single CloudEvent
to parse.
