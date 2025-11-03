package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestConvertJsonToCsv(t *testing.T) {
	jsonInput := `[
		{"id": 1, "name": "Alice", "active": true},
		{"id": 2, "name": "Bob", "active": false}
	]`

	expectedCsvOutput := "active,id,name\ntrue,1,Alice\nfalse,2,Bob\n"

	reader := strings.NewReader(jsonInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertJsonToCsv(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting JSON to CSV: %v", err)
	}

	if writer.String() != expectedCsvOutput {
		t.Errorf("Unexpected CSV output:\nExpected:\n%s\nGot:\n%s", expectedCsvOutput, writer.String())
	}
}

func TestConvertJsonToCsv_ComplexData(t *testing.T) {
	jsonInput := `[
		{"id": 1, "user": {"name": "John", "email": "john@example.com"}, "tags": ["go", "test"], "active": true, "score": null},
		{"id": 2, "user": {"name": "Jane"}, "active": false}
	]`

	expectedCsvOutput := "active,id,score,tags,user\ntrue,1,,go | test,John | john@example.com\nfalse,2,,,Jane\n"

	reader := strings.NewReader(jsonInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertJsonToCsv(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting complex JSON to CSV: %v", err)
	}

	if writer.String() != expectedCsvOutput {
		t.Errorf("Unexpected CSV output for complex data:\nExpected:\n%s\nGot:\n%s", expectedCsvOutput, writer.String())
	}
}
