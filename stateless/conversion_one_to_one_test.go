package stateless_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/test_helper"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/stretchr/testify/assert"
)

func TestConvertOneToOne(t *testing.T) {
	testcases := []struct {
		name   string
		input  message.Message[[]byte, []byte]
		output message.Message[[]byte, []byte]
		empty  bool
		err    string
	}{
		{
			name: "basic conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output: message.Message[[]byte, []byte]{
				Key:   []byte("k-updated"),
				Value: []byte("v-updated"),
			},
			err:   "",
			empty: false,
		},
		{
			name: "empty result",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			err:   "",
			empty: true,
		},
		{
			name: "error input conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("error"),
			},
			err:   "error",
			empty: false,
		},
		{
			name: "error execute",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("stupid error"),
				Value: []byte("v"),
			},
			err:   "stupid error",
			empty: false,
		},
		{
			name: "error output conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("output error"),
				Value: []byte("v"),
			},
			err:   "error",
			empty: false,
		},
	}

	crappyStringFormat := test_helper.CrappyStringFormat()
	oneToOne := func(ctx context.Context, m message.Message[string, string]) (*message.Message[string, string], error) {
		if strings.ToLower(m.Key) == "stupid error" {
			return &message.Message[string, string]{}, errors.New(m.Key)
		}

		if strings.ToLower(m.Key) == "output error" {
			return &message.Message[string, string]{
				Key:   "error",
				Value: "error",
			}, nil
		}

		if strings.ToLower(m.Key) == "empty" {
			return nil, nil
		}

		return &message.Message[string, string]{
			Key:   m.Key + "-updated",
			Value: m.Value + "-updated",
		}, nil
	}

	converted := stateless.ConvertOneToOne(oneToOne, crappyStringFormat, crappyStringFormat, crappyStringFormat, crappyStringFormat)

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			output, err := converted(context.Background(), testcase.input)
			if len(testcase.err) > 0 {
				assert.Equal(0, len(output))
				assert.Contains(err.Error(), testcase.err)
				return
			}

			if testcase.empty {
				assert.Equal(0, len(output))
				return
			}

			assert.Equal(1, len(output))
			assert.Equal(testcase.output, output[0])
		})
	}
}

func TestConvertTopicOneToOne(t *testing.T) {
	testcases := []struct {
		name   string
		input  message.Message[[]byte, []byte]
		output message.Message[[]byte, []byte]
		empty  bool
		err    string
	}{
		{
			name: "basic conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			output: message.Message[[]byte, []byte]{
				Topic: "output",
				Key:   []byte("k"),
				Value: []byte("v-updated"),
			},
			err:   "",
			empty: false,
		},
		{
			name: "empty result",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("empty"),
				Value: []byte("v"),
			},
			err:   "",
			empty: true,
		},
		{
			name: "error input conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("ghastly"),
				Value: []byte("v"),
			},
			err:   "error",
			empty: false,
		},
		{
			name: "error output conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("haunter"),
				Value: []byte("v"),
			},
			err:   "error",
			empty: false,
		},
		{
			name: "error execute",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("stupid error"),
				Value: []byte("v"),
			},
			err:   "stupid error",
			empty: false,
		},
	}

	inputTopic := topic.Generic("input", format.Gengar(), format.Gengar())
	outputTopic := topic.Generic("output", format.Gengar(), format.Gengar())

	oneToOne := func(ctx context.Context, m message.Message[string, string]) (*message.Message[string, string], error) {
		if strings.ToLower(m.Key) == "stupid error" {
			return &message.Message[string, string]{}, errors.New(m.Key)
		}

		if strings.ToLower(m.Key) == "empty" {
			return nil, nil
		}

		return &message.Message[string, string]{
			Key:   m.Key,
			Value: m.Value + "-updated",
		}, nil
	}

	converted := stateless.ConvertTopicOneToOne(oneToOne, inputTopic, outputTopic)

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			assert := assert.New(t)

			output, err := converted(context.Background(), testcase.input)
			if len(testcase.err) > 0 {
				assert.Equal(0, len(output))
				assert.Contains(err.Error(), testcase.err)
				return
			}

			if testcase.empty {
				assert.Equal(0, len(output))
				return
			}

			assert.Equal(1, len(output))
			assert.Equal(testcase.output, output[0])
		})
	}
}
