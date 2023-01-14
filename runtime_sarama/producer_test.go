package runtime_sarama_test

import (
	"testing"

	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/stretchr/testify/assert"
)

func TestProducerStartNilShouldError(t *testing.T) {
	assert := assert.New(t)

	var producer *runtime_sarama.Producer

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer is nil", err.Error())
}

func TestProducerStartNoControllerShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerBroker("localhost:12345"),
		runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
	)

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer controller is nil", err.Error())
}

func TestProducerStartEmptyBrokersShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
		runtime_sarama.WithProducerRuntimeController(runtime.NewController()),
	)

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer brokers are empty", err.Error())
}

func TestProducerStartMissingSaramaConfigShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerBroker("localhost:12345"),
		runtime_sarama.WithProducerRuntimeController(runtime.NewController()),
	)

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer sarama configuration is nil", err.Error())
}

func TestProducerStartMissingBrokerShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerBroker("localhost:12346"),
		runtime_sarama.WithProducerSaramaConfig(runtime_sarama.DefaultConfiguration()),
		runtime_sarama.WithProducerRuntimeController(runtime.NewController()),
	)

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("kafka: client has run out of available brokers to talk to: dial tcp [::1]:12346: connect: connection refused", err.Error())
}
