package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type XmlElement struct {
	XMLName  xml.Name
	Attrs    []xml.Attr   `xml:",any,attr"`
	Children []XmlElement `xml:",any"`
	Value    string       `xml:",chardata"`
}

func convertToXmlElement(data interface{}, tagName string) XmlElement {
	elem := XmlElement{XMLName: xml.Name{Local: tagName}}

	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			if array, ok := val.([]interface{}); ok {
				for _, item := range array {
					child := convertToXmlElement(item, key)
					elem.Children = append(elem.Children, child)
				}
			} else {
				child := convertToXmlElement(val, key)
				elem.Children = append(elem.Children, child)
			}

		}
	case []interface{}:
		for _, item := range v {
			elem.Children = append(elem.Children, convertToXmlElement(item, tagName))
		}
	default:
		elem.Value = fmt.Sprintf("%v", v)
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

func getJsonValue(s string) interface{} {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}

	return s
}

func writeJsonChildren(children []XmlElement, builder *strings.Builder, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	newLine := "\n"
	space := " "

	childrenGrouped := make(map[string][]XmlElement)
	var orderedKeys []string

	for _, child := range children {
		key := child.XMLName.Local
		if _, exists := childrenGrouped[key]; !exists {
			orderedKeys = append(orderedKeys, key)
		}
		childrenGrouped[key] = append(childrenGrouped[key], child)
	}

	builder.WriteString("{" + newLine)
	for i, key := range orderedKeys {
		childrenForKey := childrenGrouped[key]

		builder.WriteString(indent + "  ")
		builder.WriteString(fmt.Sprintf("%q:%s", key, space))

		if len(childrenForKey) == 1 {
			child := childrenForKey[0]

			if len(child.Children) == 0 {
				jsonValue, _ := json.Marshal(getJsonValue(child.Value))
				builder.Write(jsonValue)
			} else {
				writeJsonChildren(child.Children, builder, indentLevel+1)
			}
		} else {
			builder.WriteString("[" + newLine)
			for j, item := range childrenForKey {
				builder.WriteString(indent + "    ")
				if len(item.Children) == 0 {
					jsonValue, _ := json.Marshal(getJsonValue(item.Value))
					builder.Write(jsonValue)
				} else {
					writeJsonChildren(item.Children, builder, indentLevel+2)
				}

				if j < len(childrenForKey)-1 {
					builder.WriteString("," + newLine)
				} else {
					builder.WriteString(newLine)
				}
			}
			builder.WriteString(indent + "  " + "]")
		}

		if i < len(orderedKeys)-1 {
			builder.WriteString("," + newLine)
		} else {
			builder.WriteString(newLine)
		}
	}
	builder.WriteString(indent + "}")
}

func processXmlElement(elem XmlElement) interface{} {
	if len(elem.Children) == 0 {
		return getJsonValue(elem.Value)
	}

	childrenGrouped := make(map[string][]XmlElement)
	var orderedKeys []string

	for _, child := range elem.Children {
		key := child.XMLName.Local
		if _, exists := childrenGrouped[key]; !exists {
			orderedKeys = append(orderedKeys, key)
		}
		childrenGrouped[key] = append(childrenGrouped[key], child)
	}
	if len(orderedKeys) == 1 {
		childrenList := childrenGrouped[orderedKeys[0]]
		if len(childrenList) > 1 {
			var list []interface{}
			for _, item := range childrenList {
				list = append(list, processXmlElement(item))
			}
			return list
		}
	}
	obj := make(map[string]interface{})
	for key, childrenForKey := range childrenGrouped {
		if len(childrenForKey) == 1 {
			obj[key] = processXmlElement(childrenForKey[0])
		} else {
			var list []interface{}
			for _, item := range childrenForKey {
				list = append(list, processXmlElement(item))
			}
			obj[key] = list
		}
	}
	return obj
}

func convertJsonToCsv(input io.Reader, output io.Writer, delimiter rune) error {
	inputData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var data interface{}
	if err := json.Unmarshal(inputData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

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

func convertJsonToXml(input io.Reader, output io.Writer, rootName string) error {
	inputData, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	var data interface{}
	if err := json.Unmarshal(inputData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	xmlRoot := convertToXmlElement(data, rootName)

	xmlData, err := xml.MarshalIndent(xmlRoot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}

	xmlHeader := []byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	if _, err := output.Write(append(xmlHeader, xmlData...)); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func parseValue(s string) interface{} {
	if s == "" {
		return nil
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

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

func convertXmlToJson(input io.Reader, output io.Writer) error {

	decoder := xml.NewDecoder(input)
	var stack []*XmlElement
	var rootElement *XmlElement

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
			newElement := XmlElement{XMLName: t.Name}
			stack = append(stack, &newElement)
		case xml.EndElement:
			if len(stack) == 0 {
				return fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			current := stack[len(stack)-1]
			if current.XMLName.Local != t.Name.Local {
				return fmt.Errorf("mismatched tags: expected %s, got %s", current.XMLName.Local, t.Name.Local)
			}
			stack = stack[:len(stack)-1]

			if len(stack) == 0 {
				rootElement = current
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, *current)
			}
		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				current := stack[len(stack)-1]
				current.Value = text
			}
		}
	}

	if rootElement == nil {
		return fmt.Errorf("invalid or empty XML structure")
	}

	var builder strings.Builder
	builder.WriteString("{\n")
	builder.WriteString(fmt.Sprintf("  %q: ", rootElement.XMLName.Local))
	writeJsonChildren(rootElement.Children, &builder, 1)
	builder.WriteString("\n}")

	if _, err := output.Write([]byte(builder.String())); err != nil {
		return fmt.Errorf("failed to write JSON output file: %v", err)
	}

	return nil
}

func convertXmlToCsv(input io.Reader, output io.Writer, delimiter rune) error {

	decoder := xml.NewDecoder(input)
	var stack []*XmlElement
	var rootElement *XmlElement

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
			newElement := XmlElement{XMLName: t.Name}
			stack = append(stack, &newElement)

		case xml.EndElement:
			if len(stack) == 0 {
				return fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			current := stack[len(stack)-1]
			if current.XMLName.Local != t.Name.Local {
				return fmt.Errorf("mismatched tags: expected %s, got %s", current.XMLName.Local, t.Name.Local)
			}
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				rootElement = current
			} else {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, *current)
			}

		case xml.CharData:
			text := strings.TrimSpace(string(t))
			if text != "" && len(stack) > 0 {
				current := stack[len(stack)-1]
				current.Value = text
			}
		}
	}

	if len(stack) != 0 {
		return fmt.Errorf("invalid XML structure")
	}

	var records []interface{}
	if rootElement == nil {
		return fmt.Errorf("invalid or empty XML structure")
	}
	allSameTag := true
	if len(rootElement.Children) > 1 {
		firstTag := rootElement.Children[0].XMLName.Local
		for _, child := range rootElement.Children[1:] {
			if child.XMLName.Local != firstTag {
				allSameTag = false
				break
			}
		}
	} else {
		allSameTag = false
	}

	if allSameTag {
		for _, rowElem := range rootElement.Children {
			records = append(records, processXmlElement(rowElem))
		}
	} else {
		records = append(records, processXmlElement(*rootElement))
	}

	writer := csv.NewWriter(output)
	writer.Comma = delimiter

	if err := writeDataAsCSV(writer, records); err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return nil
}

func ensureOutputExtension(filename, desiredExt string) string {
	filename = strings.TrimSpace(filename)
	desiredExt = strings.ToLower(strings.TrimPrefix(desiredExt, "."))

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))
	if ext != desiredExt {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename))
		filename += "." + desiredExt
	}
	return filename
}
