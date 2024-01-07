package task_executor_converted

import (
	"context"

	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/structure"
)

func New[T any](
	taskChannel task.Channel[T],
	executor task.Executor[T],
) task.Executor[structure.Bytes] {
	f := &fn[T]{
		taskChannel: taskChannel,
		executor:    executor,
	}
	return f.apply
}

type fn[T any] struct {
	taskChannel task.Channel[T]
	executor    task.Executor[T]
}

func (r *fn[T]) apply(c context.Context, t task.Message[structure.Bytes]) error {
	taskConverted, taskConversionError := task.Convert(t, format.Bytes(), r.taskChannel.ValueFormat())
	if taskConversionError != nil {
		return taskConversionError
	}

	return r.executor(c, taskConverted)
}
