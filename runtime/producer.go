package runtime

import (
	"context"

	"github.com/hjwalt/flows/message"
)

type Producer interface {
	Produce(context.Context, []message.Message[message.Bytes, message.Bytes]) error
	Start() error
	Stop()
}
