package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type XmlGeneric struct {
	XmlName xml.Name
	Content []interface{} `xml:"any"`
}

func convertToXmlElement(data interface{}, tagName string) XmlGeneric {
	elem := XmlGeneric{XmlName: xml.Name{Local: tagName}}
	switch v := data.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			val := v[key]
			elem.Content = append(elem.Content, convertToXmlElement(val, key))
		}
	case []interface{}:
		for _, item := range v {
			elem.Content = append(elem.Content, convertToXmlElement(item, tagName))
		}
	case nil:

	default:
		elem.Content = append(elem.Content, fmt.Sprintf("%v", v))
	}
	return elem
}

func flattenValues(data interface{}, separator string) string {
	switch v := data.(type) {
	case map[string]interface{}:
		var parts []string
		for _, value := range v {
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

func convertJsonToCsv(inputFile, outputFile string, delimiter rune) error {
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var data interface{}
	if err := json.Unmarshal(inputData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	writer.Comma = delimiter

	if err := writeDataAsCSV(writer, data); err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}
	return nil
}

func convertJsonToXml(inputFile, outputFile string) error {
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var data interface{}
	if err := json.Unmarshal(inputData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	xmlRoot := convertToXmlElement(data, "root")
	output, err := xml.MarshalIndent(xmlRoot, "", " ")
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func convertCsvToJson(inputFile, outputFile string, delimiter rune) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	header := records[0]
	var rows []map[string]interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for j, value := range record {
			row[header[j]] = strings.TrimSpace(value)
		}
		rows = append(rows, row)
	}

	output, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write JSON output file: %v", err)
	}

	return nil
}

func convertCsvToXml(inputFile, outputFile string, delimiter rune) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %v", err)
	}

	header := records[0]
	var rows []map[string]interface{}

	for _, record := range records[1:] {
		row := make(map[string]interface{})
		for j, value := range record {
			row[header[j]] = strings.TrimSpace(value)
		}
		rows = append(rows, row)
	}

	xmlRoot := convertToXmlElement(rows, "rows")
	output, err := xml.MarshalIndent(xmlRoot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write XML output file: %v", err)
	}

	return nil
}

func convertToMap(xmlElem XmlGeneric) map[string]interface{} {
	result := make(map[string]interface{})

	if len(xmlElem.Content) == 1 {
		if str, ok := xmlElem.Content[0].(string); ok {
			result[xmlElem.XmlName.Local] = str
			return result
		}
	}

	contentMap := make(map[string]interface{})
	for _, elem := range xmlElem.Content {
		switch v := elem.(type) {
		case XmlGeneric:
			recursiveMap := convertToMap(v)
			for key, value := range recursiveMap {
				if _, exists := contentMap[key]; !exists {
					contentMap[key] = value

				} else {
					oldValue := contentMap[key]
					switch x := oldValue.(type) {
					case []interface{}:
						x = append(x, value)
						contentMap[key] = x

					default:
						contentMap[key] = []interface{}{oldValue, value}
					}
				}
			}
		}
	}
	result[xmlElem.XmlName.Local] = contentMap
	return result
}

func convertXmlToJson(inputFile, outputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}
	defer file.Close()

	var stack []XmlGeneric
	var rootElement XmlGeneric
	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to parse XML: %v", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			newElement := XmlGeneric{XmlName: t.Name}
			stack = append(stack, newElement)

		case xml.EndElement:
			if len(stack) == 0 {
				return fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			current := stack[len(stack)-1]
			if current.XmlName.Local != t.Name.Local {
				return fmt.Errorf("mismatched tags: expected %s, got %s", current.XmlName.Local, t.Name.Local)
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				rootElement = current
			} else {
				parent := &stack[len(stack)-1]
				parent.Content = append(parent.Content, current)
			}

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				current := &stack[len(stack)-1]
				current.Content = append(current.Content, text)
			}
		}
	}
	if len(stack) != 0 {
		return fmt.Errorf("invalid XML structure")
	}

	rootMap := convertToMap(rootElement)

	output, err := json.MarshalIndent(rootMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write XML output file: %v", err)
	}

	return nil
}

func convertXmlToCsv(inputFile, outputFile string, delimiter rune) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}
	defer file.Close()

	var stack []XmlGeneric
	var rootElement XmlGeneric
	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to parse XML: %v", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			newElement := XmlGeneric{XmlName: t.Name}
			stack = append(stack, newElement)

		case xml.EndElement:
			if len(stack) == 0 {
				return fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			current := stack[len(stack)-1]
			if current.XmlName.Local != t.Name.Local {
				return fmt.Errorf("mismatched tags: expected %s, got %s", current.XmlName.Local, t.Name.Local)
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				rootElement = current
			} else {
				parent := &stack[len(stack)-1]
				parent.Content = append(parent.Content, current)
			}

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				current := &stack[len(stack)-1]
				current.Content = append(current.Content, text)
			}
		}
	}
	if len(stack) != 0 {
		return fmt.Errorf("invalid XML structure")
	}

	rootMap := convertToMap(rootElement)

	csvData := interface{}(rootMap)
	for {
		currentMap, ok := csvData.(map[string]interface{})
		if !ok {
			break
		}

		if len(currentMap) != 1 {
			break
		}

		for _, value := range currentMap {
			csvData = value
		}
	}

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	writer.Comma = delimiter

	if err := writeDataAsCSV(writer, csvData); err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return nil
}
