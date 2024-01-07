package task

import (
	"context"

	"github.com/hjwalt/runway/structure"
)

type Producer interface {
	Produce(context.Context, Message[structure.Bytes]) error
	Start() error
	Stop()
}
