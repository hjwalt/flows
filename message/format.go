package message

import (
	"errors"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// for byte based storage of a message
type MessageFormat struct{}

func (helper MessageFormat) Default() Message[[]byte, []byte] {
	return Message[[]byte, []byte]{}
}

func (helper MessageFormat) Marshal(value Message[[]byte, []byte]) ([]byte, error) {
	headers := make(map[string]*protobuf.Header)
	for k, vs := range value.Headers {
		headers[k] = &protobuf.Header{
			Headers: append(make([][]byte, 0), vs...),
		}
	}

	protoMessage := &protobuf.Message{
		Topic:     value.Topic,
		Partition: value.Partition,
		Offset:    value.Offset,
		Key:       value.Key,
		Value:     value.Value,
		Headers:   headers,
		Timestamp: timestamppb.New(value.Timestamp),
	}

	protoBytes, convertErr := format.Convert(protoMessage, format.Protobuf[*protobuf.Message](), format.Bytes())
	if convertErr != nil {
		return make([]byte, 0), convertErr
	}

	return protoBytes, nil
}

func (helper MessageFormat) Unmarshal(value []byte) (Message[[]byte, []byte], error) {
	protoMessage, convertErr := format.Convert(value, format.Bytes(), format.Protobuf[*protobuf.Message]())
	if convertErr != nil {
		return helper.Default(), convertErr
	}

	headers := make(map[string][][]byte)
	for k, vs := range protoMessage.Headers {
		headers[k] = append(make([][]byte, 0), vs.Headers...)
	}

	return Message[[]byte, []byte]{
		Topic:     protoMessage.Topic,
		Partition: protoMessage.Partition,
		Offset:    protoMessage.Offset,
		Key:       protoMessage.Key,
		Value:     protoMessage.Value,
		Headers:   headers,
		Timestamp: protoMessage.Timestamp.AsTime(),
	}, nil
}

func (helper MessageFormat) ToJson(value Message[[]byte, []byte]) ([]byte, error) {
	return make([]byte, 0), errors.New("not supported")
}

func (helper MessageFormat) FromJson(value []byte) (Message[[]byte, []byte], error) {
	return helper.Default(), errors.New("not supported")
}

func Format() format.Format[Message[[]byte, []byte]] {
	return MessageFormat{}
}
