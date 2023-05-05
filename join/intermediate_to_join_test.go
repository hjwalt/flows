package join_test

import (
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/test_helper"
)

func TestIntermediateToJoinMap(t *testing.T) {
	cases := []struct {
		name   string
		key    []byte
		value  []byte
		err    error
		output []message.Message[message.Bytes, message.Bytes]
	}{
		{
			name: "normal message",
			key:  test_helper.MarshalNoError(t, message.Format(), message.Message[message.Bytes, message.Bytes]{}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
		})
	}
}
