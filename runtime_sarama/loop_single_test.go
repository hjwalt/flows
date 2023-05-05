package runtime_sarama_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/test_helper"
	"github.com/stretchr/testify/assert"
)

func TestSingleConsumeLoopWhenNoErrorShouldComplete(t *testing.T) {
	assert := assert.New(t)

	executionCount := 0

	messages := make(chan *sarama.ConsumerMessage)
	completed := make(chan bool, 1)

	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))
	consumerSingleLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(
			func(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
				executionCount += 1
				return nil, nil
			},
		),
	)

	// mock setup

	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1}), gomock.Any()).Times(1)

	// execute test

	go func() {
		err := consumerSingleLoop.Loop(session, ConsumerGroupClaimForTest{messages: messages})
		if err == nil {
			completed <- true
		} else {
			completed <- false
		}
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	messages <- &sarama.ConsumerMessage{Partition: 1}
	close(messages)
	time.Sleep(time.Millisecond)
	assert.Equal(1, len(completed))

	result := <-completed
	assert.True(result)
	assert.Equal(1, executionCount)
}

func TestSingleConsumeLoopWhenErrorShouldError(t *testing.T) {
	assert := assert.New(t)

	executionCount := 0

	messages := make(chan *sarama.ConsumerMessage)
	completed := make(chan string, 1)

	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

	consumerSingleLoop := runtime_sarama.NewSingleLoop(
		runtime_sarama.WithLoopSingleFunction(
			func(c context.Context, m message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
				executionCount += 1
				return nil, errors.New("mocked error")
			},
		),
	)

	// mock setup

	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 2}), gomock.Any()).Times(0)

	// execute test

	go func() {
		err := consumerSingleLoop.Loop(session, ConsumerGroupClaimForTest{messages: messages})
		if err == nil {
			completed <- "no error"
		} else {
			completed <- err.Error()
		}
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	messages <- &sarama.ConsumerMessage{Partition: 2}
	time.Sleep(time.Millisecond)
	assert.Equal(1, len(completed))

	result := <-completed
	assert.Equal("mocked error", result)
	assert.Equal(1, executionCount)
}
