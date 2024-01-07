package task

import (
	"context"
)

type Executor[T any] func(context.Context, Message[T]) error

type Scheduler[T any] func(context.Context) (Message[T], error)
