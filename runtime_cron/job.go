package runtime_cron

import (
	"context"

	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
)

type Job[OK any, OV any, T any] struct {
	taskProducer task.Producer
	scheduler    task.Scheduler[OK, OV, T]
	channel      task.Channel[T]
}

func (j *Job[OK, OV, T]) Run() {
	ctx := context.Background()

	t, err := j.scheduler(ctx)
	if err != nil {
		logger.ErrorErr("scheduler failed with error", err)
		return
	}

	tb, tberr := task.Convert(t, j.channel.ValueFormat(), format.Bytes())
	if tberr != nil {
		logger.ErrorErr("scheduler failed to convert task message", tberr)
		return
	}
	tb.Channel = j.channel.Name()

	taskProduceErr := j.taskProducer.Produce(ctx, tb)
	if taskProduceErr != nil {
		logger.ErrorErr("scheduler failed to publish task", taskProduceErr)
		return
	}
}
