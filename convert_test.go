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

	expectedCsvOutput := "active,id,score,tags,user\ntrue,1,,go | test,john@example.com | John\nfalse,2,,,Jane\n"

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

func TestConvertJsonToYaml(t *testing.T) {
	jsonInput := `{"id": 1, "name": "Alice", "active": true}`

	expectedYamlOutput := `
active: true
id: 1
name: "Alice"
`

	reader := strings.NewReader(jsonInput)
	writer := new(bytes.Buffer)

	err := convertJsonToYaml(reader, writer)
	if err != nil {
		t.Fatalf("Error converting JSON to YAML: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedYamlOutput) {
		t.Errorf("Unexpected YAML output:\nExpected:\n%s\nGot:\n%s", expectedYamlOutput, writer.String())
	}
}

func TestConvertJsonToYaml_ComplexData(t *testing.T) {
	jsonInput := `
{
	"user": {"name": "John", "email": "john@example.com"},
	"projects": [
		{"name": "cli-converter", "status": "done"},
		{"name": "web-server", "status": "dev"}
	],
	"tags": ["go", "yaml", "cli"]
}
`

	expectedYamlOutput := `
projects: 
  - 
    name: "cli-converter"
    status: "done"
  - 
    name: "web-server"
    status: "dev"
tags: 
  - "go"
  - "yaml"
  - "cli"
user: 
  email: "john@example.com"
  name: "John"
`

	reader := strings.NewReader(jsonInput)
	writer := new(bytes.Buffer)

	err := convertJsonToYaml(reader, writer)
	if err != nil {
		t.Fatalf("Error converting complex JSON to YAML: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedYamlOutput) {
		t.Errorf("Unexpected YAML output (complex):\nExpected:\n%s\nGot:\n%s", expectedYamlOutput, writer.String())
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

func TestConvertCsvToYaml(t *testing.T) {
	csvInput := "id,name\n1,Alice\n2,Bob"

	expectedYamlOutput := `
- 
  id: 1
  name: "Alice"
- 
  id: 2
  name: "Bob"
`
	reader := strings.NewReader(csvInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertCsvToYaml(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting CSV to YAML: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedYamlOutput) {
		t.Errorf("Unexpected YAML output (simple CSV):\nExpected:\n%s\nGot:\n%s", expectedYamlOutput, writer.String())
	}
}

func TestConvertCsvToYaml_ComplexTypes(t *testing.T) {
	csvInput := "id,is_active,score,user\n1,true,98.5,John\n2,false,100,Jane"

	expectedYamlOutput := `
- 
  id: 1
  is_active: true
  score: 98.5
  user: "John"
- 
  id: 2
  is_active: false
  score: 100
  user: "Jane"
`

	reader := strings.NewReader(csvInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertCsvToYaml(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting CSV with complex types to YAML: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedYamlOutput) {
		t.Errorf("Unexpected YAML output (complex CSV):\nExpected:\n%s\nGot:\n%s", expectedYamlOutput, writer.String())
	}
}

func TestConvertYamlToJson(t *testing.T) {
	yamlInput := `
id: 1
name: "Alice"
active: true
`

	expectedJsonOutput := `
{
  "active": true,
  "id": 1,
  "name": "Alice"
}
`
	reader := strings.NewReader(yamlInput)
	writer := new(bytes.Buffer)

	err := convertYamlToJson(reader, writer)
	if err != nil {
		t.Fatalf("Error converting YAML to JSON: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedJsonOutput) {
		t.Errorf("Unexpected JSON output (simple YAML):\nExpected:\n%s\nGot:\n%s", expectedJsonOutput, writer.String())
	}
}

func TestConvertYamlToJson_ComplexList(t *testing.T) {
	yamlInput := `
- 1782381350
- "nest"
- - "donkey"
  - general: true
    chest: true
  - - 1791717432
    - false
- 2028885901.5
`
	// O JSON de saída deve ser uma lista (array).
	expectedJsonOutput := `
[
  1782381350,
  "nest",
  [
    "donkey",
    {
      "chest": true,
      "general": true
    },
    [
      1791717432,
      false
    ]
  ],
  2028885901.5
]
`

	reader := strings.NewReader(yamlInput)
	writer := new(bytes.Buffer)

	err := convertYamlToJson(reader, writer)
	if err != nil {
		t.Fatalf("Error converting complex YAML to JSON: %v", err)
	}

	if strings.TrimSpace(writer.String()) != strings.TrimSpace(expectedJsonOutput) {
		t.Errorf("Unexpected JSON output (complex YAML):\nExpected:\n%s\nGot:\n%s", expectedJsonOutput, writer.String())
	}
}

func TestConvertYamlToCsv(t *testing.T) {
	yamlInput := `
- id: 1
  name: "Alice"
- id: 2
  name: "Bob"
`

	expectedCsvOutput := "id,name\n1,Alice\n2,Bob\n"

	reader := strings.NewReader(yamlInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertYamlToCsv(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting YAML to CSV: %v", err)
	}

	if writer.String() != expectedCsvOutput {
		t.Errorf("Unexpected CSV output (simple YAML):\nExpected:\n%s\nGot:\n%s", expectedCsvOutput, writer.String())
	}
}

func TestConvertYamlToCsv_Complex(t *testing.T) {
	yamlInput := `
- id: 1
  user:
    name: "John"
    email: "john@example.com"
  tags: ["go", "test"]
- id: 2
  user:
    name: "Jane"
  tags: ["csv"]
`

	expectedCsvOutput := "id,tags,user\n1,go | test,john@example.com | John\n2,csv,Jane\n"

	reader := strings.NewReader(yamlInput)
	writer := new(bytes.Buffer)
	delimiter := ','

	err := convertYamlToCsv(reader, writer, delimiter)
	if err != nil {
		t.Fatalf("Error converting complex YAML to CSV: %v", err)
	}

	if writer.String() != expectedCsvOutput {
		t.Errorf("Unexpected CSV output (complex YAML):\nExpected:\n%s\nGot:\n%s", expectedCsvOutput, writer.String())
	}
}
