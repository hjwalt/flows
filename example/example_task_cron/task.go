package example_task_cron

import (
	"context"

	"github.com/hjwalt/flows"
	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
)

func fn(c context.Context) (task.Message[string], error) {
	logger.Info("cron")

	return task.Message[string]{
			Value: "cron",
		},
		nil
}

func Registrar(ci inverse.Container) flows.Prebuilt {
	return flows.CronConfiguration[string]{
		Name:        Instance,
		TaskChannel: task.StringChannel("tasks"),
		Scheduler:   fn,
		Schedules: []string{
			"0 * * * * *",
			"@every 2s",
		},
		TaskConnectionString: "amqp://guest:guest@localhost:5672/",
		HttpPort:             8082,
	}
}

const (
	Instance = "tasks-example-cron"
)

func Register(m flows.Main) {
	m.Prebuilt(Instance, Registrar)
}
