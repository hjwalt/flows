package format

import (
	"encoding/json"

	reflect "github.com/hjwalt/runway/reflect"
	"gopkg.in/yaml.v3"
)

type YamlFormat[T any] struct{}

func (helper YamlFormat[T]) Default() T {
	return reflect.Construct[T]()
}

func (helper YamlFormat[T]) Marshal(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}
	return yaml.Marshal(value)
}

func (helper YamlFormat[T]) Unmarshal(value []byte) (T, error) {
	if len(value) == 0 {
		return helper.Default(), nil
	}
	yamlMessage := helper.Default()
	var err error
	if reflect.IsPointer(yamlMessage) {
		err = yaml.Unmarshal(value, yamlMessage)
	} else {
		err = yaml.Unmarshal(value, &yamlMessage)
	}
	return yamlMessage, err
}

func (helper YamlFormat[T]) ToJson(value T) ([]byte, error) {
	if reflect.IsNil(value) {
		return nil, nil
	}
	jsonbytes, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jsonbytes, err
}

func (helper YamlFormat[T]) FromJson(value []byte) (T, error) {
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

func Yaml[T any]() Format[T] {
	return YamlFormat[T]{}
}
