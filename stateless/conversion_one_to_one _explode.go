package stateless

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
)

func ConvertOneToOneExplode[IK any, IV any, OK any, OV any](
	source OneToOneExplodeFunction[IK, IV, OK, OV],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), unmarshalError
		}

		ress, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		for _, res := range ress {
			bytesResMessage, marshalError := message.Convert(res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}

func ConvertTopicOneToOneExplode[IK any, IV any, OK any, OV any](
	source OneToOneExplodeFunction[IK, IV, OK, OV],
	inputTopic topic.Topic[IK, IV],
	outputTopic topic.Topic[OK, OV],
) SingleFunction {
	return func(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
		formattedMessage, unmarshalError := message.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return make([]message.Message[[]byte, []byte], 0), unmarshalError
		}

		ress, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]message.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]message.Message[[]byte, []byte], 0)

		for _, res := range ress {
			bytesResMessage, marshalError := message.Convert(res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]message.Message[[]byte, []byte], 0), marshalError
			}
			bytesResMessage.Topic = outputTopic.Topic()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
