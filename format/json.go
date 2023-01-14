package format

import (
	"encoding/json"

	reflect "github.com/hjwalt/runway/reflect"
)

type JsonFormat[T any] struct{}

func (helper JsonFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper JsonFormat[T]) Marshal(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}
	return json.Marshal(value)
}

func (helper JsonFormat[T]) Unmarshal(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	jsonMessage := helper.Default()
	err := json.Unmarshal(value, jsonMessage)
	return jsonMessage, err
}

func (helper JsonFormat[T]) ToJson(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}
	jsonbytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jsonbytes, err
}

func (helper JsonFormat[T]) FromJson(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	jsonMessage := helper.Default()
	err := json.Unmarshal(value, jsonMessage)
	if err != nil {
		return helper.Default(), err
	}
	return jsonMessage, nil
}

func Json[T any]() Format[T] {
	return JsonFormat[T]{}
}
