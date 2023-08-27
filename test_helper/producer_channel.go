package test_helper

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/structure"
)

func NewChannelProducer() *ChannelProducer {
	return &ChannelProducer{
		Messages: make(chan flow.Message[[]byte, []byte], 100),
	}
}

type ChannelProducer struct {
	Messages chan flow.Message[structure.Bytes, structure.Bytes]
}

func (c *ChannelProducer) Produce(ctx context.Context, messages []flow.Message[structure.Bytes, structure.Bytes]) error {
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
