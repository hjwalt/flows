package runtime_sarama_test

import (
	"testing"

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

func TestProducerStartEmptyBrokersShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer()

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer brokers are empty", err.Error())
}

func TestProducerStartMissingSaramaConfigShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerBroker("localhost:12345"),
	)

	producer.(*runtime_sarama.Producer).SaramaConfiguration = nil

	err := producer.Start()

	assert.NotNil(err)
	assert.Equal("producer sarama configuration is nil", err.Error())
}

func TestProducerStartMissingBrokerShouldError(t *testing.T) {
	assert := assert.New(t)

	producer := runtime_sarama.NewProducer(
		runtime_sarama.WithProducerBroker("localhost:12346"),
	)

	err := producer.Start()

	assert.NotNil(err)
	assert.Contains(err.Error(), "kafka: client has run out of available brokers to talk to")
}
