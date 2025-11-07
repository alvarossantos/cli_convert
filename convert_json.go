package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

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

func convertJsonToYaml(input io.Reader, output io.Writer) error {
	var data interface{}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}

	return WriteAsYaml(data, output)
}
