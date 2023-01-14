package format

import (
	"encoding/json"
)

type StringFormat struct {
}

func (helper StringFormat) Default() string {
	return ""
}

func (helper StringFormat) Marshal(value string) ([]byte, error) {
	return []byte(value), nil
}

func (helper StringFormat) Unmarshal(value []byte) (string, error) {
	return string(value), nil
}

func (helper StringFormat) ToJson(value string) ([]byte, error) {
	return json.Marshal(StringJson{Message: value})
}

func (helper StringFormat) FromJson(value []byte) (string, error) {
	if len(value) == 0 {
		return "", nil
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

func String() Format[string] {
	return StringFormat{}
}
