package test_helper

import (
	"errors"
	"strings"

	"github.com/hjwalt/runway/format"
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

type StringJson struct {
	Message string
}

func CrappyStringFormat() format.Format[string] {
	return crappyStringFormat{}
}
