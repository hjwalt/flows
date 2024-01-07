package runtime_cron

import (
	"time"

	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/robfig/cron/v3"
)

// constructor
func NewCron(configurations ...runtime.Configuration[*Cron]) runtime.Runtime {
	r := &Cron{
		cron: cron.New(
			cron.WithLocation(time.UTC),
			cron.WithSeconds(),
		),
	}
	for _, configuration := range configurations {
		r = configuration(r)
	}
	return r
}

// configuration

func WithTaskProducer(taskProducer task.Producer) runtime.Configuration[*Cron] {
	return func(c *Cron) *Cron {
		c.taskProducer = taskProducer
		return c
	}
}

func WithCronJob[OK any, OV any, T any](
	schedule string,
	scheduler task.Scheduler[OK, OV, T],
	channel task.Channel[T],
) runtime.Configuration[*Cron] {
	return func(c *Cron) *Cron {
		c.cron.AddJob(
			schedule,
			&Job[OK, OV, T]{
				taskProducer: c.taskProducer,
				scheduler:    scheduler,
				channel:      channel,
			},
		)
		return c
	}
}

// implementation
type Cron struct {
	taskProducer task.Producer
	cron         *cron.Cron
}

func (c *Cron) Start() error {
	logger.Debug("cron starting")
	c.cron.Start()
	return nil
}

func (c *Cron) Stop() {
	logger.Debug("cron stopping")
	c.cron.Stop()
}
