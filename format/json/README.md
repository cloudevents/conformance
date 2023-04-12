# JSON format conformance tests

Each test file is a valid JSON document representing an object.
The top level properties are:

- `testType`: the default test type for all tests in this file;
  individual tests may override this. Values:
  - `ValidSingleEvent`
  - `ValidBatch`
  - `InvalidSingleEvent`
  - `InvalidBatch`
- `tests`: the list of tests in this file

Each test has the following properties:

- `id`: the ID of the test
- `sampleId`: for a valid event/batch test, this provides the "golden"
  sample ID to compare the result with
- `testType`: used to override the default test type in the file
- `description`: an optional description of the test (purely for clarity)
- `event`: a single CloudEvent (valid or invalid)
- `batch`: a CloudEvent batch (valid or invalid) as an array
- `sampleExtensionAttributes`: `true` if this test requires pre-defined sample extension attributes;
  implementations which do not support providing extension attribute definitions during
  parsing may skip these tests.
- `extensionConstraints`: `true` if this test requires constraints to be applied;
  implementations which do not apply constraints may skip these tests.
