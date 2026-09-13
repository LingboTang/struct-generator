package structgen

import (
	"strings"
	"testing"
)

func TestGenerate_Object(t *testing.T) {
	src, err := Generate(strings.NewReader(`{
		"user_name": "alice",
		"age": 30,
		"active": true,
		"address": {"city": "New York", "zip": "10001"},
		"scores": [9.5, 8.0]
	}`), "main")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	for _, want := range []string{
		"package main",
		"type Root struct",
		"type Address struct",
		"UserName string    `json:\"user_name\"`",
		"Scores   []float64 `json:\"scores\"`",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, src)
		}
	}
}

func TestGenerate_TopLevelArray(t *testing.T) {
	src, err := Generate(strings.NewReader(`[{"abc": "abc", "kdk": [{"kdkString": "stringValueKdkd"}]}]`), "main")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	for _, want := range []string{"type Root struct", "type KdkItem struct"} {
		if !strings.Contains(src, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, src)
		}
	}
}

func TestGenerate_CustomPackageName(t *testing.T) {
	src, err := Generate(strings.NewReader(`{"a": 1}`), "models")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if !strings.HasPrefix(src, "package models") {
		t.Errorf("expected output to start with 'package models', got:\n%s", src)
	}
}

func TestGenerate_Errors(t *testing.T) {
	cases := map[string]string{
		"empty input":          "",
		"empty array":          "[]",
		"non-object top level": `"just a string"`,
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := Generate(strings.NewReader(input), "main"); err == nil {
				t.Errorf("expected error for %s, got nil", name)
			}
		})
	}
}
