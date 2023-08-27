package stateful_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestConvertPersistenceId(t *testing.T) {
	testcases := []struct {
		name          string
		input         flow.Message[[]byte, []byte]
		err           string
		persistenceId string
	}{
		{
			name: "basic conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			err:           "",
			persistenceId: "k",
		},
		{
			name: "invalid conversion",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("v"),
			},
			err:           "error",
			persistenceId: "",
		},
		{
			name: "failed execution",
			input: flow.Message[[]byte, []byte]{
				Key:   []byte("failed"),
				Value: []byte("v"),
			},
			err:           "failed",
			persistenceId: "",
		},
	}

	crappyStringFormat := test_helper.CrappyStringFormat()
	persistenceIdFunction := func(ctx context.Context, m flow.Message[string, string]) (string, error) {
		if m.Key == "failed" {
			return "", errors.New("failed")
		}
		return m.Key, nil
	}

	converted := stateful.ConvertPersistenceId(persistenceIdFunction, crappyStringFormat, crappyStringFormat)

	for _, testcase := range testcases {

		t.Run(testcase.name, func(*testing.T) {
			assert := assert.New(t)

			persistenceId, persistenceIdErr := converted(context.Background(), testcase.input)

			if len(testcase.err) > 0 {
				assert.Equal(0, len(persistenceId))
				assert.Contains(persistenceIdErr.Error(), testcase.err)
				return
			}

			assert.Equal(testcase.persistenceId, persistenceId)
		})
	}
}
