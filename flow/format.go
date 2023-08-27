package flow

import (
	"github.com/hjwalt/flows/protobuf"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// for byte based storage of a message
type MessageFormat struct{}

func (helper MessageFormat) Default() Message[structure.Bytes, structure.Bytes] {
	return Message[structure.Bytes, structure.Bytes]{}
}

func (helper MessageFormat) Marshal(value Message[structure.Bytes, structure.Bytes]) (structure.Bytes, error) {
	headers := make(map[string]*protobuf.Header)
	for k, vs := range value.Headers {
		headers[k] = &protobuf.Header{
			Headers: append(make([]structure.Bytes, 0), vs...),
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
		return make(structure.Bytes, 0), convertErr
	}

	return protoBytes, nil
}

func (helper MessageFormat) Unmarshal(value structure.Bytes) (Message[structure.Bytes, structure.Bytes], error) {
	protoMessage, convertErr := format.Convert(value, format.Bytes(), format.Protobuf[*protobuf.Message]())
	if convertErr != nil {
		return helper.Default(), convertErr
	}

	headers := make(map[string][]structure.Bytes)
	for k, vs := range protoMessage.Headers {
		headers[k] = append(make([]structure.Bytes, 0), vs.Headers...)
	}

	return Message[structure.Bytes, structure.Bytes]{
		Topic:     protoMessage.Topic,
		Partition: protoMessage.Partition,
		Offset:    protoMessage.Offset,
		Key:       protoMessage.Key,
		Value:     protoMessage.Value,
		Headers:   headers,
		Timestamp: protoMessage.Timestamp.AsTime(),
	}, nil
}

func Format() format.Format[Message[structure.Bytes, structure.Bytes]] {
	return MessageFormat{}
}
