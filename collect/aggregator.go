package collect

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
)

func ConvertAggregator[S any, IK any, IV any](
	source Aggregator[S, IK, IV],
	s format.Format[S],
	ik format.Format[IK],
	iv format.Format[IV],
) Aggregator[message.Bytes, message.Bytes, message.Bytes] {
	return func(ctx context.Context, m message.Message[[]byte, []byte], ss stateful.State[[]byte]) (stateful.State[[]byte], error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return ss, unmarshalError
		}

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return ss, stateUnmarshalError
		}

		nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return ss, fnError
		}

		bytesNextState, stateMarshalError := stateful.ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return ss, stateMarshalError
		}

		return bytesNextState, nil
	}
}

func ConvertTopicAggregator[S any, IK any, IV any](
	source Aggregator[S, IK, IV],
	s format.Format[S],
	inputTopic topic.Topic[IK, IV],
) Aggregator[message.Bytes, message.Bytes, message.Bytes] {
	return func(ctx context.Context, m message.Message[[]byte, []byte], ss stateful.State[[]byte]) (stateful.State[[]byte], error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return ss, unmarshalError
		}

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return ss, stateUnmarshalError
		}

		nextState, fnError := source(ctx, formattedMessage, formattedState)
		if fnError != nil {
			return ss, fnError
		}

		bytesNextState, stateMarshalError := stateful.ConvertSingleState(nextState, s, format.Bytes())
		if stateMarshalError != nil {
			return ss, stateMarshalError
		}

		return bytesNextState, nil
	}
}
