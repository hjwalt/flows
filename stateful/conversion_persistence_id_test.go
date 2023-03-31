package stateful_test

import (
	"context"
	"errors"
	"testing"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/mock"
	"github.com/hjwalt/flows/stateful"
	"github.com/stretchr/testify/assert"
)

func TestConvertPersistenceId(t *testing.T) {
	testcases := []struct {
		name          string
		input         message.Message[[]byte, []byte]
		err           string
		persistenceId string
	}{
		{
			name: "basic conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("k"),
				Value: []byte("v"),
			},
			err:           "",
			persistenceId: "k",
		},
		{
			name: "invalid conversion",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("error"),
				Value: []byte("v"),
			},
			err:           "error",
			persistenceId: "",
		},
		{
			name: "failed execution",
			input: message.Message[[]byte, []byte]{
				Key:   []byte("failed"),
				Value: []byte("v"),
			},
			err:           "failed",
			persistenceId: "",
		},
	}

	crappyStringFormat := mock.CrappyStringFormat()
	persistenceIdFunction := func(ctx context.Context, m message.Message[string, string]) (string, error) {
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
				assert.EqualError(persistenceIdErr, testcase.err)
				return
			}

			assert.Equal(testcase.persistenceId, persistenceId)
		})
	}
}
