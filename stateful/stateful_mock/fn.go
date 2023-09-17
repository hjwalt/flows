package stateful_mock

import (
	"context"
	"errors"
	"strings"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
)

func MockOneToOne(ctx context.Context, m flow.Message[string, string], ss stateful.State[string]) (*flow.Message[string, string], stateful.State[string], error) {
	if strings.ToLower(m.Key) == "mock_error" {
		return nil, ss, ErrMock
	}

	if strings.ToLower(m.Key) == "recoverable" {
		return nil, ss, ErrRecoverable
	}

	if strings.ToLower(m.Key) == "irrecoverable" {
		return nil, ss, ErrIrrecoverable
	}

	ss.Content += m.Value

	if strings.ToLower(m.Key) == "empty" {
		return nil, ss, nil
	}

	if strings.ToLower(m.Key) == "mock_topic" {
		return &flow.Message[string, string]{
				Topic: "mock_topic",
				Key:   m.Key,
				Value: m.Value + "-mock_topic",
			},
			ss,
			nil
	}

	return &flow.Message[string, string]{
			Key:   m.Key,
			Value: m.Value + "-updated",
		},
		ss,
		nil
}

var (
	InputTopic       = flow.GenericTopic("input", format.Gengar(), format.Gengar())
	OutputTopic      = flow.GenericTopic("output", format.Gengar(), format.Gengar())
	EmptyTopic       = flow.GenericTopic("", format.Gengar(), format.Gengar())
	ErrMock          = errors.New("mock")
	ErrRecoverable   = errors.New("recoverable")
	ErrIrrecoverable = errors.New("irrrecoverable")
)
