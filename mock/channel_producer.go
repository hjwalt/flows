package mock

import (
	"context"

	"github.com/hjwalt/flows/message"
)

func NewChannelProducer() *ChannelProducer {
	return &ChannelProducer{
		Messages: make(chan message.Message[[]byte, []byte], 100),
	}
}

type ChannelProducer struct {
	Messages chan message.Message[message.Bytes, message.Bytes]
}

func (c *ChannelProducer) Produce(ctx context.Context, messages []message.Message[message.Bytes, message.Bytes]) error {
	for _, m := range messages {
		c.Messages <- m
	}
	return nil
}

func (c *ChannelProducer) Start() error {
	return nil
}

func (c *ChannelProducer) Stop() {

}
