package test_helper

import (
	"encoding/json"
	"errors"

	"github.com/hjwalt/flows/format"
	reflect "github.com/hjwalt/runway/reflect"
)

type CrappyJsonFormat[T CrappyShouldError] struct{}

type CrappyShouldError interface {
	ShouldError() bool
}

func (helper CrappyJsonFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper CrappyJsonFormat[T]) Marshal(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}

	if value.ShouldError() {
		return make([]byte, 0), errors.New("crappy")
	}

	return json.Marshal(value)
}

func (helper CrappyJsonFormat[T]) Unmarshal(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	jsonMessage := helper.Default()
	err := json.Unmarshal(value, jsonMessage)
	if err != nil {
		return jsonMessage, err
	}

	if jsonMessage.ShouldError() {
		return jsonMessage, errors.New("crappy")
	}

	return jsonMessage, nil
}

func (helper CrappyJsonFormat[T]) ToJson(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}

	if value.ShouldError() {
		return make([]byte, 0), errors.New("crappy")
	}

	jsonbytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jsonbytes, err
}

func (helper CrappyJsonFormat[T]) FromJson(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	jsonMessage := helper.Default()
	err := json.Unmarshal(value, jsonMessage)
	if err != nil {
		return helper.Default(), err
	}

	if jsonMessage.ShouldError() {
		return jsonMessage, errors.New("crappy")
	}

	return jsonMessage, nil
}

func CrappyJson[T CrappyShouldError]() format.Format[T] {
	return CrappyJsonFormat[T]{}
}
