package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// == test for func:JsonToString ==
func TestJsonToString_1(t *testing.T) {
	type JsonObject struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := JsonObject{
		Name: "user",
		Age:  20,
	}

	result, err := JsonToString(input)
	if err != nil {
		t.Fatalf("invalid json struct: %v", err)
	}

	expect := "{\"name\":\"user\",\"age\":20}"

	assert.Equal(t, expect, result)
}

// ================================

// == test for func:StringToJson ==
func TestStringToJson_1(t *testing.T) {
	input := "{\"name\": \"user\", \"age\": 20}"

	type JsonObject struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var result JsonObject
	if err := StringToJson(input, &result); err != nil {
		t.Fatalf("invalid json format: %v", err)
	}

	expect := JsonObject{
		Name: "user",
		Age:  20,
	}

	assert.Equal(t, expect, result)
}

// ================================

// == test for func:WriteJsonToFile ==
func TestWriteJsonToFile_1(t *testing.T) {
	t.Parallel()

	// create temporary directory
	dir := t.TempDir()

	type JsonObject struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	input := JsonObject{
		Name: "user",
		Age:  20,
	}

	path := filepath.Join(dir, "test.json")
	if err := WriteJsonToFile(path, input); err != nil {
		t.Fatalf("WriteJsonToFile failed: %v", err)
	}

	// check if the file exists
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	// validate file contents
	var result JsonObject
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid json written: %v", err)
	}

	expect := JsonObject{
		Name: "user",
		Age:  20,
	}

	assert.Equal(t, expect, result)
}

// ===================================

// == test for func:ReadJsonFile ==
func TestReadJsonFile_1(t *testing.T) {
	t.Parallel()

	// create temporary directory
	dir := t.TempDir()

	// create json file
	path := filepath.Join(dir, "test.json")
	content := `{"name": "user", "age": 20}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// test json struct
	type JsonObject struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var result JsonObject
	if err := ReadJsonFile(path, &result); err != nil {
		t.Fatalf("ReadJsonFile returned error: %v", err)
	}

	expect := JsonObject{
		Name: "user",
		Age:  20,
	}

	assert.Equal(t, expect, result)
}

// ================================
