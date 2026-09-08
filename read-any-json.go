package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"
	"unicode"
)

var inputFile string
var outputFile string

func main() {
	flag.StringVar(&inputFile, "input", "inputFile", "string of input file name")
	flag.StringVar(&outputFile, "output", "struct.go", "path to write the generated Go struct definitions to")
	flag.Parse()

	jsonFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	t, err := decoder.Token()
	if err != nil {
		if err == io.EOF {
			fmt.Printf("File %s empty.\n", inputFile)
		} else {
			fmt.Printf("Error reading token from %s: %v\n", inputFile, err)
		}
		return
	}

	if _, err := jsonFile.Seek(0, io.SeekStart); err != nil {
		fmt.Printf("Error seeking file %s: %v\n", inputFile, err)
		return
	}
	decoder = json.NewDecoder(jsonFile)

	var raw interface{}
	if err := decoder.Decode(&raw); err != nil {
		fmt.Printf("Error decoding %s: %v\n", inputFile, err)
		return
	}

	// Unwrap top-level array: use first element as representative object
	switch t.(type) {
	case json.Delim:
		delim := t.(json.Delim).String()
		if delim == "[" {
			arr, ok := raw.([]interface{})
			if !ok || len(arr) == 0 {
				fmt.Println("Empty or invalid top-level array.")
				return
			}
			raw = arr[0]
		}
	default:
		fmt.Println("Top-level JSON must be an object or array of objects.")
		return
	}

	obj, ok := raw.(map[string]interface{})
	if !ok {
		fmt.Println("Top-level JSON must be an object or array of objects.")
		return
	}

	structs := []string{}
	generateStruct("Root", obj, &structs)

	var sb strings.Builder
	sb.WriteString("package main\n\n")
	for i := len(structs) - 1; i >= 0; i-- {
		sb.WriteString(structs[i])
		sb.WriteString("\n\n")
	}

	src, err := format.Source([]byte(sb.String()))
	if err != nil {
		fmt.Printf("Error formatting generated code: %v\n", err)
		return
	}

	if err := os.WriteFile(outputFile, src, 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", outputFile, err)
		return
	}

	fmt.Printf("Wrote struct definitions to %s\n", outputFile)
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
