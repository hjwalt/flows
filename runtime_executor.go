package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_rabbit"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/flows/task_executor_converted"
	"github.com/hjwalt/flows/task_executor_retry"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

type ExecutorConfiguration[T any] struct {
	Name                        string
	TaskChannel                 task.Channel[T]
	TaskExecutor                task.Executor[T]
	TaskConnectionString        string
	HttpPort                    int
	RabbitConsumerConfiguration []runtime.Configuration[*runtime_rabbit.Consumer]
	RetryConfiguration          []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration          []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c ExecutorConfiguration[T]) Register(ci inverse.Container) {
	RegisterRabbitConsumerExecutor(
		ci,
		func(ctx context.Context, ci inverse.Container) (task.Executor[structure.Bytes], error) {
			retry, err := GetRetry(ctx, ci)
			if err != nil {
				return nil, err
			}

			executor := task_executor_converted.New[T](
				c.TaskChannel,
				c.TaskExecutor,
			)

			executor = task_executor_retry.New(
				task_executor_retry.WithRetry(retry),
				task_executor_retry.WithExecutor(executor),
			)

			return executor, nil
		},
	)

	// RUNTIME

	RegisterRetry(
		ci,
		c.RetryConfiguration,
	)
	RegisterRoute(
		ci,
		c.HttpPort,
		c.RouteConfiguration,
	)
	RegisterRabbitConsumer(
		ci,
		c.Name,
		c.TaskConnectionString,
		c.TaskChannel.Name(),
		c.RabbitConsumerConfiguration,
	)
}