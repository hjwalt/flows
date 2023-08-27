package collect

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func ConvertOneToOneCollector[S any, OK any, OV any](
	source OneToOneCollector[S, OK, OV],
	s format.Format[S],
	ok format.Format[OK],
	ov format.Format[OV],
) Collector {
	return func(ctx context.Context, persistenceId string, ss stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return flow.EmptySlice(), stateUnmarshalError
		}

		res, fnError := source(ctx, persistenceId, formattedState)
		if fnError != nil {
			return flow.EmptySlice(), fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := flow.Convert(*res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return flow.EmptySlice(), marshalError
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
	return func(ctx context.Context, persistenceId string, ss stateful.State[structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {

		formattedState, stateUnmarshalError := stateful.ConvertSingleState(ss, format.Bytes(), s)
		if stateUnmarshalError != nil {
			return flow.EmptySlice(), stateUnmarshalError
		}

		res, fnError := source(ctx, persistenceId, formattedState)
		if fnError != nil {
			return flow.EmptySlice(), fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		if res != nil {
			bytesResMessage, marshalError := flow.Convert(*res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return flow.EmptySlice(), marshalError
			}
			bytesResMessage.Topic = outputTopic.Topic()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
