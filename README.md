# struct-generator

A Go library and CLI that reads any JSON document and generates matching Go struct definitions, paired with a Python script for generating random test JSON.

- [`structgen`](./structgen) — importable Go package with the core `Generate` function.
- [`cmd/struct-generator`](./cmd/struct-generator) — CLI wrapping the package.

---

## structgen (library)

Reads a JSON document and returns generated Go struct type definitions as source code. It handles nested objects, arrays, and all standard JSON primitives, mapping them to appropriate Go types.

### Install

```bash
go get github.com/LingboTang/struct-generator/structgen
```

### Usage

```go
import "github.com/LingboTang/struct-generator/structgen"

f, _ := os.Open("config.json")
defer f.Close()

src, err := structgen.Generate(f, "models") // packageName goes into "package models"
if err != nil {
    log.Fatal(err)
}
fmt.Println(src)
```

`Generate` takes any `io.Reader` and the package name to declare the generated code under, and returns gofmt-formatted Go source as a string.

### Type Mapping

| JSON Type | Go Type |
|-----------|---------|
| string | `string` |
| number | `float64` |
| boolean | `bool` |
| object | named struct |
| array | `[]<elem type>` |
| null | `interface{}` |

Field names are converted from any casing (snake_case, camelCase, kebab-case) to PascalCase. Each field includes a `json:"..."` struct tag preserving the original key name.

For top-level arrays, the first element is used as the representative object to derive the struct shape.

---

## cmd/struct-generator (CLI)

### Prerequisites

Go installed on your machine.

### Install

```bash
go install github.com/LingboTang/struct-generator/cmd/struct-generator@latest
```

Or build/run from a clone of this repo:

```bash
go build ./cmd/struct-generator
```

```bash
go run ./cmd/struct-generator -input <path_to_json_file>
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-input` | (required) | Path to the input JSON file |
| `-output` | `struct.go` | Path to write the generated Go struct definitions to |
| `-package` | `main` | Package name for the generated file |

### Example

Input (`config.json`):

```json
{
  "user_name": "alice",
  "age": 30,
  "active": true,
  "address": {
    "city": "New York",
    "zip": "10001"
  },
  "scores": [9.5, 8.0]
}
```

Command:

```bash
go run ./cmd/struct-generator -input config.json -package models
```

Output (`struct.go`):

```go
package models

type Address struct {
	City string `json:"city"`
	Zip  string `json:"zip"`
}

type Root struct {
	UserName string    `json:"user_name"`
	Age      float64   `json:"age"`
	Active   bool      `json:"active"`
	Address  Address   `json:"address"`
	Scores   []float64 `json:"scores"`
}
```

---

## generate_json_file.py

Generates a random nested JSON file and writes it to `output.json`. Useful for testing `structgen` with varied structures.

### Prerequisites

Python 3.10+ (uses `match`/`case` syntax).

### Usage

```bash
python generate_json_file.py
```

Output is written to `output.json` in the current directory.

### How It Works

The generator builds a JSON structure recursively with these controls (configured in `main()`):

| Parameter | Description |
|-----------|-------------|
| `max_level` | Maximum nesting depth |
| `max_field_span` | Maximum number of fields per object |
| `max_array_len` | Maximum number of elements per array |
| `max_type_span` | Number of field types to include per object |

At each level, it randomly chooses to produce either an object or an array. Field types are drawn from the `FieldType` enum:

| Enum Value | Produces |
|------------|---------|
| `STRING_PAIR` | `string` value |
| `STRING_INT` | `int` value |
| `STRING_FLOAT` | `float` value |
| `STRING_DICT` | nested `{}` object |
| `STRING_LIST` | nested `[]` array |

Field keys are random 10-character alphanumeric strings.

### Typical Workflow

```bash
# 1. Generate a random JSON file
python generate_json_file.py

# 2. Generate Go structs from it
go run ./cmd/struct-generator -input output.json
```
