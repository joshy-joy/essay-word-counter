package utils

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "\n"), nil
}

func PrettyPrintJSON(data interface{}) (*bytes.Buffer, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	return &out, err
}
