package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func convertCsvToJson(input io.Reader, output io.Writer, delimiter rune) error {
	file, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(file)))
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("empty CSV file")
	}

	header := records[0]
	var rows []map[string]interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for j, value := range record {
			cleanValue := strings.TrimSpace(value)
			row[header[j]] = parseValue(cleanValue)
		}
		rows = append(rows, row)
	}

	jsonBytes, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	_, err = output.Write(jsonBytes)
	return err
}

func convertCsvToXml(input io.Reader, output io.Writer, delimiter rune, rootName string) error {
	reader := csv.NewReader(input)
	reader.Comma = delimiter

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(records) < 1 {
		return fmt.Errorf("empty CSV file")
	}

	header := records[0]
	var rows []map[string]interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for j, value := range record {

			cleanValue := strings.TrimSpace(value)
			if num, err := strconv.ParseFloat(cleanValue, 64); err == nil {
				row[header[j]] = num
			} else {
				row[header[j]] = cleanValue
			}
		}
		rows = append(rows, row)
	}

	xmlRoot := XmlElement{XMLName: xml.Name{Local: rootName}}
	for _, rowMap := range rows {
		rowElem := convertToXmlElement(rowMap, "row")
		xmlRoot.Children = append(xmlRoot.Children, rowElem)
	}

	xmlData, err := xml.MarshalIndent(xmlRoot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %v", err)
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	if _, err := output.Write(append(xmlHeader, xmlData...)); err != nil {
		return fmt.Errorf("failed to write XML output file: %v", err)
	}

	return nil
}

func convertCsvToYaml(input io.Reader, output io.Writer, delimiter rune) error {
	reader := csv.NewReader(input)
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("empty CSV file")
	}

	header := records[0]
	var rows []interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for j, value := range record {
			cleanValue := strings.TrimSpace(value)
			row[header[j]] = parseValue(cleanValue)
		}
		rows = append(rows, row)
	}
	return WriteAsYaml(rows, output)
}

func flattenValues(data interface{}, separator string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		var parts []string
		for _, key := range keys {
			value := v[key]
			parts = append(parts, flattenValues(value, separator))
		}
		return strings.Join(parts, separator)
	case []interface{}:
		var parts []string
		for _, item := range v {
			parts = append(parts, flattenValues(item, separator))
		}
		return strings.Join(parts, separator)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func writeDataAsCSV(writer *csv.Writer, data interface{}) error {
	var rows []interface{}

	switch v := data.(type) {
	case []interface{}:
		if len(v) == 0 {
			return nil
		}
		rows = v
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("empty object")
		}
		rows = []interface{}{v}
	default:
		return fmt.Errorf("format not supported")
	}

	headerSet := make(map[string]struct{})
	for _, row := range rows {
		if obj, ok := row.(map[string]interface{}); ok {
			for key := range obj {
				headerSet[key] = struct{}{}
			}
		}
	}

	headers := make([]string, 0, len(headerSet))
	for key := range headerSet {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, row := range rows {
		if obj, ok := row.(map[string]interface{}); ok {
			record := make([]string, len(headers))
			for i, header := range headers {
				if value, exists := obj[header]; exists {
					switch v := value.(type) {
					case map[string]interface{}, []interface{}:
						record[i] = flattenValues(v, " | ")
					case nil:
						record[i] = ""
					default:
						record[i] = fmt.Sprintf("%v", v)
					}
				} else {
					record[i] = ""
				}
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	return nil
}
