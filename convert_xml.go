package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func convertXmlToJson(input io.Reader, output io.Writer) error {

	rootElement, err := parseXmlToElement(input)
	if err != nil {
		return err
	}

	if rootElement == nil {
		return fmt.Errorf("invalid or empty XML structure")
	}

	data := processXmlElement(*rootElement)

	result := map[string]interface{}{
		rootElement.XMLName.Local: data,
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	if _, err := output.Write(jsonBytes); err != nil {
		return fmt.Errorf("failed to write JSON output file: %v", err)
	}

	return nil
}

func convertXmlToCsv(input io.Reader, output io.Writer, delimiter rune) error {

	rootElement, err := parseXmlToElement(input)
	if err != nil {
		return err
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

	err = writeDataAsCSV(writer, records)
	if err != nil {
		return err
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %v", err)
	}

	return nil
}

func convertXmlToYaml(input io.Reader, output io.Writer) error {
	rootElement, err := parseXmlToElement(input)
	if err != nil {
		return err
	}

	data := processXmlElement(*rootElement)

	result := map[string]interface{}{
		rootElement.XMLName.Local: data,
	}

	if err := WriteAsYaml(result, output); err != nil {
		return err
	}

	return nil
}

type XmlElement struct {
	XMLName  xml.Name
	Attrs    []xml.Attr   `xml:",any,attr"`
	Children []XmlElement `xml:",any"`
	Value    string       `xml:",chardata"`
}

func parseXmlToElement(input io.Reader) (*XmlElement, error) {
	decoder := xml.NewDecoder(input)
	var stack []*XmlElement
	var rootElement *XmlElement

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse XML: %v", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			newElement := XmlElement{XMLName: t.Name}
			stack = append(stack, &newElement)
		case xml.EndElement:
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected end element %s", t.Name.Local)
			}
			current := stack[len(stack)-1]
			if current.XMLName.Local != t.Name.Local {
				return nil, fmt.Errorf("mismatched tags: expected %s, got %s", current.XMLName.Local, t.Name.Local)
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
		return nil, fmt.Errorf("invalid or empty XML structure")
	}
	return rootElement, nil
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

func getJsonValue(s string) interface{} {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}

	return s
}
