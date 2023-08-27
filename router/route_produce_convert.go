package router

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func RouteProduceBodyMapConvert[Req any, Key any, Value any](
	source stateless.OneToOneFunction[structure.Bytes, Req, Key, Value],
	requestFormat format.Format[Req],
	keyFormat format.Format[Key],
	valueFormat format.Format[Value],
) stateless.OneToOneFunction[structure.Bytes, structure.Bytes, structure.Bytes, structure.Bytes] {
	return func(ctx context.Context, m flow.Message[[]byte, []byte]) (*flow.Message[[]byte, []byte], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), format.Bytes(), requestFormat)
		if unmarshalError != nil {
			return nil, unmarshalError
		}

		outputMessage, outputerr := source(ctx, formattedMessage)
		if outputerr != nil {
			return nil, outputerr
		}

		if outputMessage == nil {
			return nil, nil
		}

		bytesResMessage, marshalError := flow.Convert(*outputMessage, keyFormat, valueFormat, format.Bytes(), format.Bytes())
		if marshalError != nil {
			return nil, marshalError
		}

		return &bytesResMessage, nil
	}
}

func RouteProduceTopicBodyMapConvert[Req any, Key any, Value any](
	source stateless.OneToOneFunction[structure.Bytes, Req, Key, Value],
	requestFormat format.Format[Req],
	produceTopic topic.Topic[Key, Value],
) stateless.OneToOneFunction[structure.Bytes, structure.Bytes, structure.Bytes, structure.Bytes] {
	return func(ctx context.Context, m flow.Message[[]byte, []byte]) (*flow.Message[[]byte, []byte], error) {

		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), format.Bytes(), requestFormat)
		if unmarshalError != nil {
			return nil, unmarshalError
		}

		outputMessage, outputerr := source(ctx, formattedMessage)
		if outputerr != nil {
			return nil, outputerr
		}

		if outputMessage == nil {
			return nil, nil
		}

		bytesResMessage, marshalError := flow.Convert(*outputMessage, produceTopic.KeyFormat(), produceTopic.ValueFormat(), format.Bytes(), format.Bytes())
		if marshalError != nil {
			return nil, marshalError
		}
		bytesResMessage.Topic = produceTopic.Topic()

		return &bytesResMessage, nil
	}
}
