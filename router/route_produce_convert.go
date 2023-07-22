package router

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
)

func RouteProduceBodyMapConvert[Req any, Key any, Value any](
	source stateless.OneToOneFunction[message.Bytes, Req, Key, Value],
	requestFormat format.Format[Req],
	keyFormat format.Format[Key],
	valueFormat format.Format[Value],
) stateless.OneToOneFunction[message.Bytes, message.Bytes, message.Bytes, message.Bytes] {
	return func(ctx context.Context, m message.Message[[]byte, []byte]) (*message.Message[[]byte, []byte], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), format.Bytes(), requestFormat)
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

		bytesResMessage, marshalError := message.Convert(*outputMessage, keyFormat, valueFormat, format.Bytes(), format.Bytes())
		if marshalError != nil {
			return nil, marshalError
		}

		return &bytesResMessage, nil
	}
}

func RouteProduceTopicBodyMapConvert[Req any, Key any, Value any](
	source stateless.OneToOneFunction[message.Bytes, Req, Key, Value],
	requestFormat format.Format[Req],
	produceTopic topic.Topic[Key, Value],
) stateless.OneToOneFunction[message.Bytes, message.Bytes, message.Bytes, message.Bytes] {
	return func(ctx context.Context, m message.Message[[]byte, []byte]) (*message.Message[[]byte, []byte], error) {

		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), format.Bytes(), requestFormat)
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

		bytesResMessage, marshalError := message.Convert(*outputMessage, produceTopic.KeyFormat(), produceTopic.ValueFormat(), format.Bytes(), format.Bytes())
		if marshalError != nil {
			return nil, marshalError
		}
		bytesResMessage.Topic = produceTopic.Topic()

		return &bytesResMessage, nil
	}
}
