package main

import (
	"fmt"
	"encoding/json"
	"os"
	"io"
	"flag"
)

var (
	inputFile string
)

func main() {

	flag.StringVar(&inputFile, "input", "inputFile", "string of input file name")
	//flag.StringVar(&outputFile, "output", "outputFile", "string of output file name")

	flag.Parse()

	fmt.Printf("Input file: %s\n", inputFile)
	//fmt.Println("Output file: ", outputFile)

	jsonFile, err := os.Open(inputFile)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened json file")


	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)

	t, err := decoder.Token()
	if err != nil {
		if err == io.EOF {
			fmt.Printf("File %s empty.\n\n", jsonFile)
		} else {
			fmt.Printf("Error peeking token from %s: %v\n\n", jsonFile, err)
		}
		return
	}

	fmt.Printf("  First JSON token type: %T, Value: %v\n", t, t)

	if _, err := jsonFile.Seek(0, io.SeekStart); err != nil {
		fmt.Printf("Error seeking file %s: %v\n\n", inputFile, err)
	}

	decoder = json.NewDecoder(jsonFile)


	//var data map[string]interface{}

	switch v := t.(type) {
		case json.Delim: // Indicates a JSON object '{' or array '['
			if v.String() == "{" {
				fmt.Println("  Detected: JSON Object")
				var data map[string]interface{} // Declare for object
				if err := decoder.Decode(&data); err != nil {
					fmt.Printf("  Error decoding object from %s: %v\n", jsonFile, err)
					return
				}
				fmt.Println("  Decoded object (unstructured, field by field):")
				for key, val := range data {
					fmt.Printf("    Field: %s, Value: %v, GoType: %T\n", key, val, val)
					// Further type-specific processing can happen here
				}
			} else if v.String() == "[" {
				fmt.Println("  Detected: JSON Array")
				var data []interface{} // Declare for array
				if err := decoder.Decode(&data); err != nil {
					fmt.Printf("  Error decoding array from %s: %v\n", jsonFile, err)
					return
				}
				fmt.Println("  Decoded array (unstructured, element by element):")
				for i, elem := range data {
					fmt.Printf("    Element[%d]: %v, GoType: %T\n", i, elem, elem)
					// If array elements are objects, you can cast them:
					if objElem, ok := elem.(map[string]interface{}); ok {
						fmt.Println("      (Array element is an object, iterating fields):")
						for objKey, objVal := range objElem {
							fmt.Printf("        Sub-Field: %s, Value: %v, GoType: %T\n", objKey, objVal, objVal)
						}
					}
				}
			}
		case bool:
			fmt.Println("  Detected: JSON Boolean")
			var data bool
			if err := decoder.Decode(&data); err != nil {
				fmt.Printf("  Error decoding boolean from %s: %v\n", jsonFile, err)
				return
			}
			fmt.Printf("  Decoded boolean value: %t\n", data)
		case float64: // JSON numbers
			fmt.Println("  Detected: JSON Number")
			var data float64
			if err := decoder.Decode(&data); err != nil {
				fmt.Printf("  Error decoding number from %s: %v\n", jsonFile, err)
				return
			}
			fmt.Printf("  Decoded number value: %f\n", data)
		case string:
			fmt.Println("  Detected: JSON String")
			var data string
			if err := decoder.Decode(&data); err != nil {
				fmt.Printf("  Error decoding string from %s: %v\n", jsonFile, err)
				return
			}
			fmt.Printf("  Decoded string value: %q\n", data)
		case nil: // JSON null
			fmt.Println("  Detected: JSON Null")
			var data interface{} // null decodes to nil Go interface{}
			if err := decoder.Decode(&data); err != nil {
				fmt.Printf("  Error decoding null from %s: %v\n", jsonFile, err)
				return
			}
			fmt.Printf("  Decoded null value: %v\n", data)
		default:
			fmt.Printf("  Detected: Unknown top-level JSON type: %T\n", v)
		}
		fmt.Println("-----------------------------------------\n")


	


}