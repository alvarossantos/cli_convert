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

	// As chaves do JSON (active, id, name) são ordenadas alfabeticamente para gerar o cabeçalho do CSV.
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

	// As chaves (active, id, score, tags, user) são ordenadas alfabeticamente.
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

func TestConvertCsvToJson(t *testing.T) {
	csvInput := "id,name,age\n1,John,30\n2,Jane,25\n"

	expectedJsonOutput := `[
  {
    "age": 30,
    "id": 1,
    "name": "John"
  },
  {
    "age": 25,
    "id": 2,
    "name": "Jane"
  }
]`

	input := strings.NewReader(csvInput)
	var output bytes.Buffer
	delimiter := ','

	err := convertCsvToJson(input, &output, delimiter)
	if err != nil {
		t.Fatalf("Error converting CSV to JSON: %v", err)
	}

	if strings.TrimSpace(output.String()) != strings.TrimSpace(expectedJsonOutput) {
		t.Errorf("Unexpected JSON output:\nExpected:\n%s\nGot:\n%s", expectedJsonOutput, output.String())
	}
}

func TestConvertCsvToJson_WithDifferentTypes(t *testing.T) {
	csvInput := "id,name,is_active,score\n1,John,true,98.5\n2,Jane,false,100\n3,João,true,85\n"

	expectedJsonOutput := `[
  {
    "id": 1,
    "is_active": true,
    "name": "John",
    "score": 98.5
  },
  {
    "id": 2,
    "is_active": false,
    "name": "Jane",
    "score": 100
  },
  {
    "id": 3,
    "is_active": true,
    "name": "João",
    "score": 85
  }
]`

	input := strings.NewReader(csvInput)
	var output bytes.Buffer
	delimiter := ','

	err := convertCsvToJson(input, &output, delimiter)
	if err != nil {
		t.Fatalf("Error converting CSV with types to JSON: %v", err)
	}

	if strings.TrimSpace(output.String()) != strings.TrimSpace(expectedJsonOutput) {
		t.Errorf("Unexpected JSON output for typed CSV:\nExpected:\n%s\nGot:\n%s", expectedJsonOutput, output.String())
	}
}
