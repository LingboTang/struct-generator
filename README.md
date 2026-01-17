# struct-generator
Generic struct generator based on unstructured json and other input files


Read Any JSON
A flexible Go utility that reads a generic JSON file, detects its top-level type (Object, Array, or Primitive), and prints the structure and values to the console.

Prerequisites
Go (Golang) installed on your machine.

Installation
Save the code as read-any-json.go.

Open your terminal or command prompt.

Build the executable:

Bash
go build read-any-json.go
This will create a binary file named read-any-json (or read-any-json.exe on Windows).

Usage
You can run the tool by pointing it to a JSON file using the -input flag.

Basic Command

Bash
./read-any-json -input <path_to_json_file>
(Note: On Windows, use read-any-json.exe -input <path_to_json_file>)

Using go run directly

If you do not want to build a binary, you can run the source code directly:

Bash
go run read-any-json.go -input data.json
Examples
1. Reading a JSON Object

Input (config.json):

JSON
{
  "host": "localhost",
  "port": 8080
}
Command:

Bash
./read-any-json -input config.json
Output:

Plaintext
Input file: config.json
Successfully Opened json file
  First JSON token type: json.Delim, Value: {
  Detected: JSON Object
  Decoded object (unstructured, field by field):
    Field: host, Value: localhost, GoType: string
    Field: port, Value: 8080, GoType: float64
2. Reading a JSON Array

Input (users.json):

JSON
["alice", "bob", "charlie"]
Command:

Bash
./read-any-json -input users.json
Output:

Plaintext
  Detected: JSON Array
  Decoded array (unstructured, element by element):
    Element[0]: alice, GoType: string
    Element[1]: bob, GoType: string
    ...
How It Works
Flag Parsing: The program accepts an input filename via the command line flag.

Token Peeking: It reads the very first token of the file to determine if the JSON is an Object ({), Array ([), or a primitive value (String, Number, Boolean, Null).

Rewinding: It rewinds the file pointer to the beginning.

Dynamic Decoding: Based on the detected type, it creates the appropriate Go data structure (e.g., map[string]interface{} or []interface{}) and decodes the file.

Inspection: It iterates over the decoded data and prints the keys, values, and underlying Go types.