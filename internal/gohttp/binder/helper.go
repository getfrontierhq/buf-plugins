package binder

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// parseProtobufTag extracts the field name and delimiter from the struct tag
func parseTag(tag string) string {
	parts := strings.Split(tag, structTagConfigDelimiter)
	fieldName := ""
	for _, part := range parts {
		if strings.HasPrefix(part, "name=") {
			fieldName = strings.TrimPrefix(part, "name=")
			continue
		}
		if strings.HasPrefix(part, "json=") {
			fieldName = strings.TrimPrefix(part, "json=")
			break
		}
	}

	return fieldName
}

// setFieldValue sets the struct field based on its type, parsing the query value accordingly
func setFieldValue(fieldVal reflect.Value, queryValue, delimiter string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(queryValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseAndSetInt(fieldVal, queryValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseAndSetUint(fieldVal, queryValue)
	case reflect.Float32, reflect.Float64:
		return parseAndSetFloat(fieldVal, queryValue)
	case reflect.Bool:
		return parseAndSetBool(fieldVal, queryValue)
	case reflect.Slice:
		return parseAndSetSlice(fieldVal, queryValue, delimiter)
	default:
		return fmt.Errorf("unsupported field type: %v", fieldVal.Kind())
	}
	return nil
}

// parseAndSetInt parses and sets an integer value
func parseAndSetInt(fieldVal reflect.Value, queryValue string) error {
	intValue, err := strconv.ParseInt(queryValue, 10, fieldVal.Type().Bits())
	if err != nil {
		return fmt.Errorf("parsing int: %v", err)
	}
	fieldVal.SetInt(intValue)
	return nil
}

// parseAndSetUint parses and sets an unsigned integer value
func parseAndSetUint(fieldVal reflect.Value, queryValue string) error {
	uintValue, err := strconv.ParseUint(queryValue, 10, fieldVal.Type().Bits())
	if err != nil {
		return fmt.Errorf("parsing uint: %v", err)
	}
	fieldVal.SetUint(uintValue)
	return nil
}

// parseAndSetFloat parses and sets a floating-point value
func parseAndSetFloat(fieldVal reflect.Value, queryValue string) error {
	floatValue, err := strconv.ParseFloat(queryValue, fieldVal.Type().Bits())
	if err != nil {
		return fmt.Errorf("parsing float: %v", err)
	}
	fieldVal.SetFloat(floatValue)
	return nil
}

// parseAndSetBool parses and sets a boolean value
func parseAndSetBool(fieldVal reflect.Value, queryValue string) error {
	boolValue, err := strconv.ParseBool(queryValue)
	if err != nil {
		return fmt.Errorf("parsing bool: %v", err)
	}
	fieldVal.SetBool(boolValue)
	return nil
}

// parseAndSetSlice parses and sets a slice value, splitting the query value based on the specified delimiter
func parseAndSetSlice(fieldVal reflect.Value, queryValue, delimiter string) error {
	parts := strings.Split(queryValue, delimiter)
	slice := reflect.MakeSlice(fieldVal.Type(), len(parts), len(parts))

	for i, part := range parts {
		item := slice.Index(i)

		// Special handling for []interface{}
		if item.Kind() == reflect.Interface {
			parsedValue, err := parseToInterface(part)
			if err != nil {
				return err
			}

			item.Set(reflect.ValueOf(parsedValue))
			continue
		}

		if err := setFieldValue(item, part, delimiter); err != nil {
			return err
		}
	}

	fieldVal.Set(slice)
	return nil
}

// parseToInterface attempts to parse a string to the most appropriate basic type for interface{} usage
func parseToInterface(value string) (interface{}, error) {
	// Try parsing as int
	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return intValue, nil
	}
	// Try parsing as float
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue, nil
	}
	// Try parsing as bool
	if boolValue, err := strconv.ParseBool(value); err == nil {
		return boolValue, nil
	}
	// Default to string
	return value, nil
}

func decodeURL(pattern, path string) (map[string]string, error) {
	splittedPattern := strings.Split(strings.Trim(pattern, "/"), "/")
	splittedPath := strings.Split(strings.Trim(path, "/"), "/")

	if len(splittedPattern) != len(splittedPath) {
		return nil, fmt.Errorf("parsing url: expect path, %s; got path, %s", pattern, path)
	}

	pathParams := make(map[string]string)
	for i, patternPart := range splittedPattern {
		if !strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			continue
		}

		key := patternPart[1 : len(patternPart)-1]
		val := splittedPath[i]

		pathParams[key] = val
	}

	return pathParams, nil
}

func extractPathVars(path string) map[string]*string {
	splittedPath := strings.Split(strings.Trim(path, "/"), "/")

	pathVars := make(map[string]*string)
	for i, pathPart := range splittedPath {
		if !strings.HasPrefix(pathPart, "{") && strings.HasSuffix(pathPart, "}") {
			continue
		}

		key := pathPart[1 : len(pathPart)-1]
		pathVars[key] = &splittedPath[i]
	}

	return pathVars
}

func encodeURL(pattern string, vars map[string]*string) string {
	splittedPattern := strings.Split(strings.Trim(pattern, "/"), "/")

	for i, patternPart := range splittedPattern {
		if !strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			continue
		}

		key := patternPart[1 : len(patternPart)-1]
		val, ok := vars[key]
		if !ok || val == nil {
			splittedPattern[i] = ""
			continue
		}

		splittedPattern[i] = *val
	}

	return strings.Join(splittedPattern, "/")
}

func hasBody(r *http.Request) bool {
	if r.Body == nil {
		return false
	}

	// Read up to 1 byte from the body
	buf := make([]byte, 1)
	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	// If we read 1 byte, there’s a body
	if n > 0 {
		// Put the byte back for later processing
		r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buf[:n]), r.Body))
		return true
	}

	return false
}

func shouldHaveBody(method string) bool {
	return method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete
}
