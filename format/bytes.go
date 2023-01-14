package format

import (
	"encoding/base64"
	"encoding/json"
)

type BytesFormat struct {
}

func (helper BytesFormat) Default() []byte {
	return []byte{}
}

func (helper BytesFormat) Marshal(value []byte) ([]byte, error) {
	return value, nil
}

func (helper BytesFormat) Unmarshal(value []byte) ([]byte, error) {
	return value, nil
}

func (helper BytesFormat) ToJson(value []byte) ([]byte, error) {
	return json.Marshal(BytesJson{Message: base64.StdEncoding.EncodeToString(value)})
}

func (helper BytesFormat) FromJson(value []byte) ([]byte, error) {
	if len(value) == 0 {
		return []byte{}, nil
	}
	jsonMessage := &BytesJson{}
	err := json.Unmarshal(value, jsonMessage)
	if err != nil {
		return []byte{}, err
	}
	return base64.StdEncoding.DecodeString(jsonMessage.Message)
}

type BytesJson struct {
	Message string
}

func Bytes() Format[[]byte] {
	return BytesFormat{}
}
