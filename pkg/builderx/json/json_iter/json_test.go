package json_iter

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pathTestData = "data"

func init() {
	if _, err := os.Stat(pathTestData); os.IsNotExist(err) {
		os.Mkdir(pathTestData, 0755)
	}
}

type Adf struct {
	Version  float64   `json:"_v"`
	Sections []Section `json:"sections"`
}

type Section struct {
	ID     string  `json:"_id"`
	Tittle string  `json:"tittle"`
	Fields []Field `json:"fields"`
}

type Field struct {
	ID         string `json:"_id"`
	Key        string `json:"key"`
	Type       string `json:"type"`
	Directive  string `json:"directive"`
	IsRequired bool   `json:"isRequired"`
}

func TestUnmarshal(t *testing.T) {
	testData := []byte(`{"name":"John","age":30,"city":"New York"}`)
	var result map[string]interface{}
	err := Unmarshal(testData, &result)
	assert.NoError(t, err)
	assert.Equal(t, "John", result["name"])
	assert.Equal(t, float64(30), result["age"])
	assert.Equal(t, "New York", result["city"])
}

func TestMarshal(t *testing.T) {
	data := map[string]interface{}{
		"name": "John",
		"age":  30,
		"city": "New York",
	}
	expected := []byte(`{"name":"John","age":30,"city":"New York"}`)
	result, err := Marshal(data)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), string(result))
}
