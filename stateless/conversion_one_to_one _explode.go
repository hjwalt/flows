package stateless

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func ConvertOneToOneExplode[IK any, IV any, OK any, OV any](
	source OneToOneExplodeFunction[IK, IV, OK, OV],
	ik format.Format[IK],
	iv format.Format[IV],
	ok format.Format[OK],
	ov format.Format[OV],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), ik, iv)
		if unmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), unmarshalError
		}

		ress, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]flow.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		for _, res := range ress {
			bytesResMessage, marshalError := flow.Convert(res, ok, ov, format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), marshalError
			}
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}

func ConvertTopicOneToOneExplode[IK any, IV any, OK any, OV any](
	source OneToOneExplodeFunction[IK, IV, OK, OV],
	inputTopic flow.Topic[IK, IV],
	outputTopic flow.Topic[OK, OV],
) SingleFunction {
	return func(ctx context.Context, m flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
		formattedMessage, unmarshalError := flow.Convert(m, format.Bytes(), format.Bytes(), inputTopic.KeyFormat(), inputTopic.ValueFormat())
		if unmarshalError != nil {
			return make([]flow.Message[[]byte, []byte], 0), unmarshalError
		}

		ress, fnError := source(ctx, formattedMessage)
		if fnError != nil {
			return make([]flow.Message[[]byte, []byte], 0), fnError
		}

		byteResultMessages := make([]flow.Message[[]byte, []byte], 0)

		for _, res := range ress {
			bytesResMessage, marshalError := flow.Convert(res, outputTopic.KeyFormat(), outputTopic.ValueFormat(), format.Bytes(), format.Bytes())
			if marshalError != nil {
				return make([]flow.Message[[]byte, []byte], 0), marshalError
			}
			bytesResMessage.Topic = outputTopic.Name()
			byteResultMessages = append(byteResultMessages, bytesResMessage)
		}

		return byteResultMessages, nil
	}
}
