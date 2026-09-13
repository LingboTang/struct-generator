// Package structgen generates Go struct definitions from JSON documents.
package structgen

import (
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"strings"
	"unicode"
)

// Generate reads a JSON document from r and returns generated Go struct
// definitions covering its shape, declared under the given package name.
//
// The top-level JSON value must be an object, or an array of objects — in
// which case the first element is used as the representative shape. Nested
// objects and arrays are recursively turned into their own named structs.
func Generate(r io.Reader, packageName string) (string, error) {
	var raw interface{}
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("input is empty")
		}
		return "", fmt.Errorf("decoding JSON: %w", err)
	}

	if arr, ok := raw.([]interface{}); ok {
		if len(arr) == 0 {
			return "", fmt.Errorf("top-level array is empty")
		}
		raw = arr[0]
	}

	obj, ok := raw.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("top-level JSON must be an object or array of objects")
	}

	structs := []string{}
	generateStruct("Root", obj, &structs)

	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", packageName)
	for i := len(structs) - 1; i >= 0; i-- {
		sb.WriteString(structs[i])
		sb.WriteString("\n\n")
	}

	src, err := format.Source([]byte(sb.String()))
	if err != nil {
		return "", fmt.Errorf("formatting generated code: %w", err)
	}
	return string(src), nil
}

// generateStruct emits a Go struct definition for a JSON object and recurses
// into nested objects/arrays. Results are appended to structs in depth-first
// order so callers can reverse for top-down output.
func generateStruct(name string, obj map[string]interface{}, structs *[]string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "type %s struct {\n", name)

	for key, val := range obj {
		fieldName := toPascalCase(key)
		goType := resolveType(fieldName, val, structs)
		fmt.Fprintf(&sb, "\t%s %s `json:\"%s\"`\n", fieldName, goType, key)
	}

	sb.WriteString("}")
	*structs = append(*structs, sb.String())
	return name
}

// resolveType returns the Go type string for a JSON value, recursing into
// nested objects and arrays and registering new structs as needed.
func resolveType(fieldName string, val interface{}, structs *[]string) string {
	switch v := val.(type) {
	case map[string]interface{}:
		generateStruct(fieldName, v, structs)
		return fieldName
	case []interface{}:
		elemType := resolveArrayElemType(fieldName, v, structs)
		return "[]" + elemType
	case string:
		return "string"
	case float64:
		return "float64"
	case bool:
		return "bool"
	case nil:
		return "interface{}"
	default:
		return "interface{}"
	}
}

// resolveArrayElemType inspects the first element of a JSON array to determine
// the element Go type.
func resolveArrayElemType(fieldName string, arr []interface{}, structs *[]string) string {
	if len(arr) == 0 {
		return "interface{}"
	}
	return resolveType(fieldName+"Item", arr[0], structs)
}

// toPascalCase converts a JSON key (snake_case, camelCase, kebab-case) to PascalCase.
func toPascalCase(s string) string {
	var sb strings.Builder
	capitalizeNext := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			capitalizeNext = true
			continue
		}
		if capitalizeNext {
			sb.WriteRune(unicode.ToUpper(r))
			capitalizeNext = false
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
