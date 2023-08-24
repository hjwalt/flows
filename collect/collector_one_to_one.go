package collect

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
)

func ConvertOneToOneCollector[S any, OK any, OV any](
	source OneToOneCollector[S, OK, OV],
	s format.Format[S],
	ok format.Format[OK],
	ov format.Format[OV],
) Collector {
	return func(ctx context.Context, persistenceId string, ss stateful.State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return message.EmptySlice(), stateUnmarshalError
		}

		res, fnError := source(ctx, persistenceId, formattedState)
		if fnError != nil {
			return message.EmptySlice(), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return message.EmptySlice(), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}

func ConvertTopicOneToOneCollector[S any, OK any, OV any](
	source OneToOneCollector[S, OK, OV],
	s format.Format[S],
	outputTopic topic.Topic[OK, OV],
) Collector {
	return func(ctx context.Context, persistenceId string, ss stateful.State[message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return message.EmptySlice(), stateUnmarshalError
		}

		res, fnError := source(ctx, persistenceId, formattedState)
		if fnError != nil {
			return message.EmptySlice(), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := message.Convert(*res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return message.EmptySlice(), marshalError
			}
			bytesResMessage.Topic = outputTopic.Topic()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
