# struct-generator

A Go tool that reads any JSON file and generates Go struct definitions from it, paired with a Python script for generating random test JSON.

---

## read-any-json.go

Reads a JSON file and prints Go struct type definitions to stdout. It handles nested objects, arrays, and all standard JSON primitives, mapping them to appropriate Go types.

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

### Prerequisites

Go installed on your machine.

### Build

```bash
go build read-any-json.go
```

### Usage

```bash
./read-any-json -input <path_to_json_file>
```

Or run directly without building:

```bash
go run read-any-json.go -input <path_to_json_file>
```

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
go run read-any-json.go -input config.json
```

Output:

```go
type Address struct {
    City string `json:"city"`
    Zip  string `json:"zip"`
}

type Root struct {
    UserName string  `json:"user_name"`
    Age      float64 `json:"age"`
    Active   bool    `json:"active"`
    Address  Address `json:"address"`
    Scores   []float64 `json:"scores"`
}
```

---

## generate_json_file.py

Generates a random nested JSON file and writes it to `output.json`. Useful for testing `read-any-json.go` with varied structures.

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
go run read-any-json.go -input output.json
```
