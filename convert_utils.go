package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var ErrLossyConversion = errors.New("data was flattened; hierarchical structure is lost")

func parseValue(s string) interface{} {
	s = strings.TrimSpace(s)

	if (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, `'`) && strings.HasSuffix(s, `'`)) {
		if len(s) > 1 {
			return s[1 : len(s)-1]
		}
		return ""
	}
	if s == "" || s == "null" {
		return nil
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") ||
		strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {

		var jsonData interface{}
		if err := json.Unmarshal([]byte(s), &jsonData); err == nil {
			return jsonData
		}
	}

	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
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

func WriteAsYaml(data interface{}, output io.Writer) error {
	return writeYamlRecursive(output, data, 0)
}

func writeYamlRecursive(writer io.Writer, data interface{}, indentLevel int) error {
	var err error

	switch v := data.(type) {
	case map[string]interface{}:
		indent := strings.Repeat("  ", indentLevel)

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, key := range keys {
			value := v[key]

			if i > 0 {
				_, err = fmt.Fprint(writer, "\n")
				if err != nil {
					return err
				}
			}

			_, err = fmt.Fprintf(writer, "%s%s: ", indent, key)
			if err != nil {
				return err
			}

			_, isMap := value.(map[string]interface{})
			_, isSlice := value.([]interface{})

			if isMap || isSlice {
				_, err = fmt.Fprint(writer, "\n")
				if err != nil {
					return err
				}

				err := writeYamlRecursive(writer, value, indentLevel+1)
				if err != nil {
					return err
				}
			} else {
				err = writeYamlRecursive(writer, value, indentLevel)
				if err != nil {
					return err
				}
			}
		}

	case []interface{}:
		indent := strings.Repeat("  ", indentLevel)

		for i, item := range v {
			if i > 0 {
				_, err = fmt.Fprint(writer, "\n")
				if err != nil {
					return err
				}
			}

			_, err = fmt.Fprintf(writer, "%s- ", indent)
			if err != nil {
				return err
			}

			_, isMap := item.(map[string]interface{})
			_, isSlice := item.([]interface{})

			if isMap || isSlice {
				_, err = fmt.Fprint(writer, "\n")
				if err != nil {
					return err
				}

				err = writeYamlRecursive(writer, item, indentLevel+1)
				if err != nil {
					return err
				}
			} else {
				err = writeYamlRecursive(writer, item, indentLevel)
				if err != nil {
					return err
				}
			}
		}

	default:
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal simple value: %v", err)
		}

		_, err = writer.Write(jsonValue)
		if err != nil {
			return err
		}
	}
	return nil
}
