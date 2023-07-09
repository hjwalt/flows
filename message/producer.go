package message

import (
	"context"
)

type Producer interface {
	Produce(context.Context, []Message[Bytes, Bytes]) error
	Start() error
	Stop()
}
