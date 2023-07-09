package runtime_sarama_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/stretchr/testify/assert"
)

// Test utilities

type ConsumerLoopForTest struct {
	loopError bool
}

func (c ConsumerLoopForTest) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	if c.loopError {
		return errors.New("mocked loop error")
	}
	return nil
}
func (batchConsume ConsumerLoopForTest) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (batchConsume ConsumerLoopForTest) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (c ConsumerLoopForTest) Start() error {
	return nil
}
func (c ConsumerLoopForTest) Stop() {
}

type ConsumerGroupForTest struct {
	consumeError bool
}

func (cgt ConsumerGroupForTest) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	if cgt.consumeError {
		return errors.New("mocked consume error")
	}
	return nil
}
func (cgt ConsumerGroupForTest) Errors() <-chan error {
	return make(<-chan error)
}
func (cgt ConsumerGroupForTest) Close() error {
	return nil
}
func (cgt ConsumerGroupForTest) Pause(partitions map[string][]int32)  {}
func (cgt ConsumerGroupForTest) Resume(partitions map[string][]int32) {}
func (cgt ConsumerGroupForTest) PauseAll()                            {}
func (cgt ConsumerGroupForTest) ResumeAll()                           {}

// Test codes

func TestConsumerStartNilShouldError(t *testing.T) {
	assert := assert.New(t)

	var consumer *runtime_sarama.Consumer

	err := consumer.Start()

	assert.ErrorIs(err, runtime_sarama.ErrConsumerIsNil)
}

func TestConsumerStartEmptyTopicsShouldError(t *testing.T) {
	assert := assert.New(t)

	consumer := runtime_sarama.NewConsumer(
		runtime_sarama.WithConsumerBroker("test-broker:9092"),
		runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
		runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
	)

	err := consumer.Start()

	assert.ErrorIs(err, runtime_sarama.ErrConsumerTopicsEmpty)
}

func TestConsumerStartEmptyBrokersShouldError(t *testing.T) {
	assert := assert.New(t)

	consumer := runtime_sarama.NewConsumer(
		runtime_sarama.WithConsumerTopic("test-topic"),
		runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
		runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
	)

	err := consumer.Start()

	assert.ErrorIs(err, runtime_sarama.ErrConsumerBrokersEmpty)
}

func TestConsumerStartEmptyConsumerLoopShouldError(t *testing.T) {
	assert := assert.New(t)

	consumer := runtime_sarama.NewConsumer(
		runtime_sarama.WithConsumerTopic("test-topic"),
		runtime_sarama.WithConsumerBroker("test-broker:9092"),
		runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
	)

	err := consumer.Start()

	assert.ErrorIs(err, runtime_sarama.ErrConsumerLoopIsNil)
}

func TestConsumerStartMissingSaramaConfigShouldError(t *testing.T) {
	assert := assert.New(t)

	consumer := runtime_sarama.NewConsumer(
		runtime_sarama.WithConsumerTopic("test-topic"),
		runtime_sarama.WithConsumerBroker("test-broker:9092"),
		runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
	)

	err := consumer.Start()

	assert.ErrorIs(err, runtime_sarama.ErrConsumerSaramaConfigurationIsNil)
}

func TestConsumerStartMissingBrokerShouldError(t *testing.T) {
	assert := assert.New(t)

	consumer := runtime_sarama.NewConsumer(
		runtime_sarama.WithConsumerBroker("localhost:12345"),
		runtime_sarama.WithConsumerTopic("test-topic"),
		runtime_sarama.WithConsumerSaramaConfig(runtime_sarama.DefaultConfiguration()),
		runtime_sarama.WithConsumerLoop(ConsumerLoopForTest{}),
	)

	err := consumer.Start()

	assert.NotNil(err)
	assert.Contains(err.Error(), "kafka: client has run out of available brokers to talk to")
}

func TestConsumerRunWhenNoConsumeErrorShouldExitAfterCancel(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())

	completed := make(chan bool, 1)

	consumer := runtime_sarama.Consumer{
		Context: ctx,
		Cancel:  cancel,
		Group:   ConsumerGroupForTest{consumeError: false},
	}

	go func() {
		consumer.Run()
		completed <- true
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	cancel()
	time.Sleep(time.Millisecond)
	assert.Equal(1, len(completed))
}

func TestConsumerRunWhenConsumeErrorShouldExitImmediately(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())

	completed := make(chan bool, 1)

	consumer := runtime_sarama.Consumer{
		Context: ctx,
		Cancel:  cancel,
		Group:   ConsumerGroupForTest{consumeError: true},
	}

	go func() {
		consumer.Run()
		completed <- true
	}()
	time.Sleep(time.Millisecond)
	assert.Equal(1, len(completed))
}
