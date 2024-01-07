package task

import (
	"context"
)

type Executor[T any] func(context.Context, Message[T]) error

type Scheduler[K any, V any, T any] func(context.Context) (Message[T], error)
