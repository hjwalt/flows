package format

import (
	"encoding/base64"
	"encoding/json"
)

type Base64Format struct {
}

func (helper Base64Format) Default() string {
	return ""
}

func (helper Base64Format) Marshal(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(value)
}

func (helper Base64Format) Unmarshal(value []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(value), nil
}

func (helper Base64Format) ToJson(value string) ([]byte, error) {
	return json.Marshal(StringJson{Message: value})
}

func (helper Base64Format) FromJson(value []byte) (string, error) {
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

func Base64() Format[string] {
	return Base64Format{}
}
