// Command struct-generator reads a JSON file and writes generated Go struct
// definitions covering its shape.
package main

import (
	"fmt"
	"os"

	"flag"

	"github.com/LingboTang/struct-generator/structgen"
)

func main() {
	var inputFile, outputFile, packageName string
	flag.StringVar(&inputFile, "input", "", "path to the input JSON file (required)")
	flag.StringVar(&outputFile, "output", "struct.go", "path to write the generated Go struct definitions to")
	flag.StringVar(&packageName, "package", "main", "package name for the generated file")
	flag.Parse()

	if inputFile == "" {
		fmt.Println("missing required -input flag")
		os.Exit(1)
	}

	jsonFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	src, err := structgen.Generate(jsonFile, packageName)
	if err != nil {
		fmt.Printf("Error generating structs from %s: %v\n", inputFile, err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, []byte(src), 0644); err != nil {
		fmt.Printf("Error writing %s: %v\n", outputFile, err)
		os.Exit(1)
	}

	fmt.Printf("Wrote struct definitions to %s\n", outputFile)
}
