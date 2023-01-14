package utility

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/hjwalt/runway/logger"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

type TestContainerLogging struct {
}

func (TestContainerLogging) Printf(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func LogContains(t *testing.T, container testcontainers.Container, logString string) bool {
	containerLogs, logError := container.Logs(context.Background())
	assert.Nil(t, logError)

	if b, err := io.ReadAll(containerLogs); err == nil {
		return strings.Contains(string(b), logString)
	}
	return false
}

func Log(t *testing.T, container testcontainers.Container) {
	containerLogs, logError := container.Logs(context.Background())
	assert.Nil(t, logError)

	if b, err := io.ReadAll(containerLogs); err == nil {
		logger.Info(string(b))
	}
}
