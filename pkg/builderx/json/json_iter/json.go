package json_iter

import (
	"bytes"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Get(data []byte, key string) (size int, res interface{}, err error) {
	val := json.Get(data, key)
	if val.LastError() != nil {
		return 0, nil, val.LastError()
	}

	return val.Size(), val.ToString(), nil
}

func SetOrdered(path string, newPath string) error {
	stream, err := os.Open(path)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.NewDecoder(stream).Decode(&data); err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}

	if err := os.WriteFile(newPath, buf.Bytes(), 0644); err != nil {
		return err
	}

	// do somet clearing garbage
	defer func() {
		err = stream.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	return nil
}

// Generate io.Reader from string
func GenerateReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
