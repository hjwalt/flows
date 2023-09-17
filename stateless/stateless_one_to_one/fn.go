package stateless_one_to_one

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func New[IK any, IV any, OK any, OV any](
	source stateless.OneToOneFunction[IK, IV, OK, OV],
	inputTopic flow.Topic[IK, IV],
	outputTopic flow.Topic[OK, OV],
) stateless.SingleFunction {
	res := &fn[IK, IV, OK, OV]{
		source:      source,
		inputKey:    inputTopic.KeyFormat(),
		inputValue:  inputTopic.ValueFormat(),
		outputKey:   outputTopic.KeyFormat(),
		outputValue: outputTopic.ValueFormat(),
		outputTopic: outputTopic.Name(),
	}
	return res.apply
}

type fn[IK any, IV any, OK any, OV any] struct {
	source      stateless.OneToOneFunction[IK, IV, OK, OV]
	inputKey    format.Format[IK]
	inputValue  format.Format[IV]
	outputKey   format.Format[OK]
	outputValue format.Format[OV]
	outputTopic string
}

func (r *fn[IK, IV, OK, OV]) apply(c context.Context, m flow.Message[structure.Bytes, structure.Bytes]) ([]flow.Message[structure.Bytes, structure.Bytes], error) {
	formattedMessage, unmarshalError := flow.Convert(
		m,
		format.Bytes(),
		format.Bytes(),
		r.inputKey,
		r.inputValue,
	)
	if unmarshalError != nil {
		return make([]flow.Message[[]byte, []byte], 0), unmarshalError
	}

	res, fnError := r.source(c, formattedMessage)
	if fnError != nil {
		return make([]flow.Message[[]byte, []byte], 0), fnError
	}

	if res == nil {
		return flow.EmptySlice(), nil
	}

	bytesResMessage, marshalError := flow.Convert(
		*res,
		r.outputKey,
		r.outputValue,
		format.Bytes(),
		format.Bytes(),
	)
	if marshalError != nil {
		return make([]flow.Message[[]byte, []byte], 0), marshalError
	}

	if len(r.outputTopic) > 0 {
		bytesResMessage.Topic = r.outputTopic
	}

	return []flow.Message[[]byte, []byte]{bytesResMessage}, nil
}
