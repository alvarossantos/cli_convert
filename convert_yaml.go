package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

func convertYamlToJson(input io.Reader, output io.Writer) error {
	data, err := parseYamlToInterface(input)
	if err != nil {
		return err
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	_, err = output.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}

func convertYamlToCsv(input io.Reader, output io.Writer, delimiter rune) error {
	data, err := parseYamlToInterface(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	writer := csv.NewWriter(output)
	writer.Comma = delimiter

	err = writeDataAsCSV(writer, data)
	if err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return nil
}

func convertYamlToXml(input io.Reader, output io.Writer, rootName string) error {
	data, err := parseYamlToInterface(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
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

type stackItem struct {
	indent    int
	container interface{}
}

func normalizeStructure(value interface{}) interface{} {
	switch v := value.(type) {
	case *map[string]interface{}:
		return normalizeMap(*v)
	case map[string]interface{}:
		return normalizeMap(v)

	case *[]interface{}:
		return normalizeSlice(*v)
	case []interface{}:
		return normalizeSlice(v)
	default:
		return v
	}
}

func normalizeMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(m))
	for key, value := range m {
		result[key] = normalizeStructure(value)
	}
	return result
}

func normalizeSlice(s []interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, value := range s {
		result[i] = normalizeStructure(value)
	}
	return result
}

func getIndent(line string) (int, string) {
	indent := 0
	for i, ch := range line {
		if ch == ' ' {
			indent++
		} else {
			return indent, line[i:]
		}
	}
	return indent, ""
}

func parseYamlToInterface(input io.Reader) (interface{}, error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var root interface{}
	var stack []stackItem
	var firstLineFound bool

	for _, line := range lines {
		_, trimmedLine := getIndent(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		if strings.HasPrefix(trimmedLine, "- ") {
			rootSlice := make([]interface{}, 0)
			root = &rootSlice
			stack = []stackItem{{indent: -1, container: &rootSlice}}
		} else {
			rootMap := make(map[string]interface{})
			root = &rootMap
			stack = []stackItem{{indent: -1, container: &rootMap}}
		}
		firstLineFound = true
		break
	}

	if !firstLineFound {
		rootMap := make(map[string]interface{})
		root = &rootMap
		stack = []stackItem{{indent: -1, container: &rootMap}}
	}

	for i, line := range lines {
		indent, trimmedLine := getIndent(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		for len(stack) > 1 && indent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}

		parent := stack[len(stack)-1].container

		if strings.HasPrefix(trimmedLine, "- ") {
			parentSlice, ok := parent.(*[]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid indentation at line %d: parent map not found", i+1)
			}

			itemStr := strings.TrimSpace(trimmedLine[2:])

			if itemStr == "" {
				isList := false
				if i+1 < len(lines) {
					_, nextLineTrimmed := getIndent(lines[i+1])
					if strings.HasPrefix(nextLineTrimmed, "- ") {
						isList = true
					}
				}

				if isList {
					newSlice := make([]interface{}, 0)
					*parentSlice = append(*parentSlice, &newSlice)
					stack = append(stack, stackItem{indent: indent, container: &newSlice})
				} else {
					newMap := make(map[string]interface{})
					*parentSlice = append(*parentSlice, &newMap)
					stack = append(stack, stackItem{indent: indent, container: &newMap})
				}
			} else if strings.HasPrefix(itemStr, "- ") {
				newSlice := make([]interface{}, 0)
				nestedItemStr := strings.TrimSpace(itemStr[2:])
				newSlice = append(newSlice, parseValue(nestedItemStr))
				*parentSlice = append(*parentSlice, &newSlice)
				stack = append(stack, stackItem{indent: indent + 1, container: &newSlice})

			} else if strings.Contains(itemStr, ":") {
				parts := strings.SplitN(itemStr, ":", 2)
				key := strings.TrimSpace(parts[0])
				valueStr := strings.TrimSpace(parts[1])

				newMap := make(map[string]interface{})
				if valueStr != "" {
					newMap[key] = parseValue(valueStr)
				}

				*parentSlice = append(*parentSlice, &newMap)
				stack = append(stack, stackItem{indent: indent + 1, container: &newMap})
			} else {
				*parentSlice = append(*parentSlice, parseValue(itemStr))
			}
			continue
		}

		if strings.Contains(trimmedLine, ":") {
			parts := strings.SplitN(trimmedLine, ":", 2)
			key := strings.TrimSpace(parts[0])
			valueStr := strings.TrimSpace(parts[1])

			parentMap, ok := parent.(*map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid yaml at line %d: found key/value pair but parent is not a map", i+1)
			}

			if valueStr == "" {
				isList := false
				if i+1 < len(lines) {
					_, nextLineTrimmed := getIndent(lines[i+1])
					if strings.HasPrefix(nextLineTrimmed, "- ") {
						isList = true
					}
				}

				if isList {
					newSlice := make([]interface{}, 0)
					(*parentMap)[key] = &newSlice
					stack = append(stack, stackItem{indent: indent, container: &newSlice})
				} else {
					newMap := make(map[string]interface{})
					(*parentMap)[key] = &newMap
					stack = append(stack, stackItem{indent: indent, container: &newMap})
				}
			} else {
				(*parentMap)[key] = parseValue(valueStr)
			}
			continue
		}
		return nil, fmt.Errorf("invalid yaml syntax at line %d: %s", i+1, trimmedLine)

	}
	cleaned := normalizeStructure(root)
	return cleaned, nil
}
