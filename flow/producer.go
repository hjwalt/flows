package flow

import (
	"context"

	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

type Producer interface {
	Produce(context.Context, []Message[structure.Bytes, structure.Bytes]) error
	Start() error
	Stop()
}

func Produce[K any, V any](
	c context.Context,
	flowProducer Producer,
	topic Topic[K, V],
	m Message[K, V],
) error {
	byteMessage, convertErr := Convert(m, topic.KeyFormat(), topic.ValueFormat(), format.Bytes(), format.Bytes())
	if convertErr != nil {
		return convertErr
	}
	byteMessage.Topic = topic.Name()
	return flowProducer.Produce(c, []Message[structure.Bytes, structure.Bytes]{byteMessage})
}
