package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonToString_Success(t *testing.T) {
	// == arrange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	input := JsonStruct{Name: "user", Age: 20}

	// == act ==
	got, err := JsonToString(input)

	// == assert ==
	assert.Equal(t, `{"name":"user","age":20}`, got)
	assert.Nil(t, err)
}

func TestJsonToString_InvalidJsonStructError(t *testing.T) {
	// == arrange ==
	input := make(chan int)

	// == act ==
	got, err := JsonToString(input)

	// == assert ==
	assert.Equal(t, "", got)
	assert.NotNil(t, err)
}

func TestStringToJson_Success(t *testing.T) {
	// == arrange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var jsonObject JsonStruct
	input := `{"name":"user","age":20}`

	// == act ==
	err := StringToJson(input, &jsonObject)

	// == assert ==
	assert.Equal(t, JsonStruct{Name: "user", Age: 20}, jsonObject)
	assert.Nil(t, err)
}

func TestStrtingToJson_InvalidJsonFormatError(t *testing.T) {
	// == arrange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var jsonObject JsonStruct
	input := "Non Json-Formatted String"

	// == act ==
	err := StringToJson(input, &jsonObject)

	// == assert ==
	assert.NotNil(t, err)
}

func TestWriteJsonToFile_Success(t *testing.T) {
	// == arrange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	jsonObject := JsonStruct{Name: "user", Age: 20}
	path := filepath.Join(t.TempDir(), "test.json")

	// == act ==
	err := WriteJsonToFile(path, jsonObject)

	// == assert ==
	assert.Nil(t, err)

	// read created file
	content, fileOpenErr := os.ReadFile(path)
	if fileOpenErr != nil {
		t.Fatalf("file open failed")
	}
	assert.Equal(t, `{
    "name": "user",
    "age": 20
}
`, string(content))
}

func TestWriteJsonToFile_PathNotExistsError(t *testing.T) {
	// == arrange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	jsonObject := JsonStruct{Name: "user", Age: 20}
	path := "/not/exists/path"

	// == act ==
	err := WriteJsonToFile(path, jsonObject)

	// == assert ==
	assert.NotNil(t, err)
}

func TestReadJsonFile_Success(t *testing.T) {
	// arrange
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	jsonObject := JsonStruct{Name: "user", Age: 20}
	path := filepath.Join(t.TempDir(), "test.json")

	// create test json file
	f, createErr := os.Create(path)
	if createErr != nil {
		t.Fatalf("create json file failed")
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.Encode(jsonObject)

	// == act ==
	var got JsonStruct
	err := ReadJsonFile(path, &got)

	// == assert ==
	assert.Nil(t, err)
	assert.Equal(t, jsonObject, got)
}

func TestReadJsonFile_FileNotExistsError(t *testing.T) {
	// == arange ==
	type JsonStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	path := "/not/exists/path"

	// == act ==
	var got JsonStruct
	err := ReadJsonFile(path, &got)

	// == assert ==
	assert.NotNil(t, err)
}
