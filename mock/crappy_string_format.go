package mock

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/hjwalt/flows/format"
)

type crappyStringFormat struct {
}

func (helper crappyStringFormat) Default() string {
	return ""
}

func (helper crappyStringFormat) Marshal(value string) ([]byte, error) {
	if strings.ToLower(value) == "error" {
		return []byte{}, errors.New(value)
	}
	return []byte(value), nil
}

func (helper crappyStringFormat) Unmarshal(value []byte) (string, error) {
	if strings.ToLower(string(value)) == "error" {
		return "", errors.New(string(value))
	}
	return string(value), nil
}

func (helper crappyStringFormat) ToJson(value string) ([]byte, error) {
	if strings.ToLower(value) == "error" {
		return []byte{}, errors.New(value)
	}
	return json.Marshal(StringJson{Message: value})
}

func (helper crappyStringFormat) FromJson(value []byte) (string, error) {
	if len(value) == 0 {
		return "", nil
	}
	if strings.ToLower(string(value)) == "error" {
		return "", errors.New(string(value))
	}
	jsonMessage := &StringJson{}
	err := json.Unmarshal(value, jsonMessage)
	if err != nil {
		return "", err
	}
	return jsonMessage.Message, nil
}

type StringJson struct {
	Message string
}

func CrappyStringFormat() format.Format[string] {
	return crappyStringFormat{}
}
