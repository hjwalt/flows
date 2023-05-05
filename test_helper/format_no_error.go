package test_helper

import (
	"testing"

	"github.com/hjwalt/flows/format"
	"github.com/stretchr/testify/assert"
)

func MarshalNoError[T any](t *testing.T, f format.Format[T], v T) []byte {
	bytes, err := f.Marshal(v)
	assert.NoError(t, err)
	return bytes
}

func UnmarshalNoError[T any](t *testing.T, f format.Format[T], bytes []byte) T {
	v, err := f.Unmarshal(bytes)
	assert.NoError(t, err)
	return v
}
