package runtime_sarama_test

import (
	"context"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/test_helper"
	"github.com/hjwalt/runway/logger"
	"github.com/stretchr/testify/assert"
)

type ConsumerGroupClaimForTest struct {
	messages chan *sarama.ConsumerMessage
}

func (cgct ConsumerGroupClaimForTest) Topic() string {
	return ""
}
func (cgct ConsumerGroupClaimForTest) Partition() int32 {
	return 0
}
func (cgct ConsumerGroupClaimForTest) InitialOffset() int64 {
	return 0
}
func (cgct ConsumerGroupClaimForTest) HighWaterMarkOffset() int64 {
	return 0
}
func (cgct ConsumerGroupClaimForTest) Messages() <-chan *sarama.ConsumerMessage {
	return cgct.messages
}

func TestKeyedHandlerShouldTriggerOnMaxBuffered(t *testing.T) {
	assert := assert.New(t)

	executionCount := 0

	messages := make(chan *sarama.ConsumerMessage)
	completed := make(chan bool, 1)

	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

	consumerBatchLoop := runtime_sarama.NewKeyedHandler(
		runtime_sarama.WithKeyedHandlerMaxBufferred(2),
		runtime_sarama.WithKeyedHandlerMaxDelay(100*time.Millisecond),
		runtime_sarama.WithKeyedHandlerFunction(
			func(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
				executionCount += 1
				assert.Equal(2, len(m))
				return nil, nil
			},
		),
		runtime_sarama.WithKeyedHandlerKeyFunction(func(ctx context.Context, m message.Message[[]byte, []byte]) (string, error) {
			return string(m.Key), nil
		}),
	)

	// mock setup

	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1, Offset: 2, Key: []byte("test2")}), gomock.Any()).Times(1)

	// execute test

	go func() {
		err := consumerBatchLoop.ConsumeClaim(session, ConsumerGroupClaimForTest{messages: messages})
		if err == nil {
			completed <- true
		} else {
			completed <- false
		}
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 1, Key: []byte("test1")}
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 2, Key: []byte("test2")}
	time.Sleep(time.Millisecond)
	close(messages)

	var result bool

	select {
	case result = <-completed:
		logger.Info("completed")
	case <-time.After(1 * time.Second):
		result = false
	}

	assert.True(result)
	assert.Equal(1, executionCount)
}

func TestKeyedHandlerShouldTriggerOnSameKey(t *testing.T) {
	assert := assert.New(t)

	executionCount := 0

	messages := make(chan *sarama.ConsumerMessage)
	completed := make(chan bool, 1)

	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

	consumerBatchLoop := runtime_sarama.NewKeyedHandler(
		runtime_sarama.WithKeyedHandlerMaxBufferred(100),
		runtime_sarama.WithKeyedHandlerMaxDelay(10*time.Millisecond),
		runtime_sarama.WithKeyedHandlerFunction(
			func(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
				executionCount += 1
				if executionCount == 1 {
					assert.Equal(3, len(m))
				}
				return nil, nil
			},
		),
		runtime_sarama.WithKeyedHandlerKeyFunction(func(ctx context.Context, m message.Message[[]byte, []byte]) (string, error) {
			return string(m.Key), nil
		}),
	)

	// mock setup

	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1, Offset: 3, Key: []byte("test3")}), gomock.Any()).Times(1)
	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1, Offset: 4, Key: []byte("test1")}), gomock.Any()).Times(1)

	// execute test

	go func() {
		err := consumerBatchLoop.ConsumeClaim(session, ConsumerGroupClaimForTest{messages: messages})
		if err == nil {
			completed <- true
		} else {
			completed <- false
		}
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 1, Key: []byte("test1")}
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 2, Key: []byte("test2")}
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 3, Key: []byte("test3")}
	messages <- &sarama.ConsumerMessage{Partition: 1, Offset: 4, Key: []byte("test1")}
	time.Sleep(11 * time.Millisecond)
	close(messages)

	var result bool

	select {
	case result = <-completed:
		logger.Info("completed")
	case <-time.After(1 * time.Second):
		result = false
	}

	assert.True(result)
	assert.Equal(2, executionCount)
}

func TestKeyedHandlerWhenNoErrorShouldTriggerOnTimer(t *testing.T) {
	assert := assert.New(t)

	executionCount := 0

	messages := make(chan *sarama.ConsumerMessage)
	completed := make(chan bool, 1)

	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

	consumerBatchLoop := runtime_sarama.NewKeyedHandler(
		runtime_sarama.WithKeyedHandlerMaxBufferred(2),
		runtime_sarama.WithKeyedHandlerMaxDelay(10*time.Millisecond),
		runtime_sarama.WithKeyedHandlerFunction(
			func(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
				executionCount += 1
				assert.Equal(1, len(m))
				return nil, nil
			},
		),
		runtime_sarama.WithKeyedHandlerKeyFunction(func(ctx context.Context, m message.Message[[]byte, []byte]) (string, error) {
			return string(m.Key), nil
		}),
	)

	// mock setup

	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1, Key: []byte("test")}), gomock.Any()).Times(1)

	// execute test

	go func() {
		err := consumerBatchLoop.ConsumeClaim(session, ConsumerGroupClaimForTest{messages: messages})
		if err == nil {
			completed <- true
		} else {
			completed <- false
		}
	}()

	time.Sleep(time.Millisecond)
	assert.Equal(0, len(completed))
	messages <- &sarama.ConsumerMessage{Partition: 1, Key: []byte("test")}
	time.Sleep(20 * time.Millisecond)
	close(messages)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(1, len(completed))

	result := <-completed
	assert.True(result)
	assert.Equal(1, executionCount)
}

// func TestBatchConsumeLoopWhenErrorShouldErrorOnMaxBuffered(t *testing.T) {
// 	assert := assert.New(t)

// 	executionCount := 0

// 	messages := make(chan *sarama.ConsumerMessage)
// 	completed := make(chan bool, 1)

// 	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

// 	consumerBatchLoop := runtime_sarama.NewBatchLoop(
// 		runtime_sarama.WithLoopBatchMaxBufferred(2),
// 		runtime_sarama.WithLoopBatchMaxDelay(100*time.Millisecond),
// 		runtime_sarama.WithLoopBatchFunction(
// 			func(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
// 				executionCount += 1
// 				assert.Equal(2, len(m))
// 				return nil, errors.New("mocked error")
// 			},
// 		),
// 	)

// 	// mock setup

// 	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1}), gomock.Any()).Times(0)

// 	// execute test

// 	go func() {
// 		err := consumerBatchLoop.ConsumeClaim(session, ConsumerGroupClaimForTest{messages: messages})
// 		assert.NotNil(err)
// 		assert.Equal("mocked error", err.Error())
// 		if err == nil {
// 			completed <- true
// 		} else {
// 			completed <- false
// 		}
// 	}()

// 	time.Sleep(time.Millisecond)
// 	assert.Equal(0, len(completed))
// 	messages <- &sarama.ConsumerMessage{Partition: 1}
// 	messages <- &sarama.ConsumerMessage{Partition: 1}
// 	close(messages)
// 	time.Sleep(101 * time.Millisecond)
// 	assert.Equal(1, len(completed))

// 	result := <-completed
// 	assert.False(result)
// 	assert.Equal(1, executionCount)
// }

// func TestBatchConsumeLoopWhenErrorShouldErrorOnTimer(t *testing.T) {
// 	assert := assert.New(t)

// 	executionCount := 0

// 	messages := make(chan *sarama.ConsumerMessage)
// 	completed := make(chan bool, 1)

// 	session := test_helper.NewMockConsumerGroupSession(gomock.NewController(t))

// 	consumerBatchLoop := runtime_sarama.NewBatchLoop(
// 		runtime_sarama.WithLoopBatchMaxBufferred(2),
// 		runtime_sarama.WithLoopBatchMaxDelay(100*time.Millisecond),
// 		runtime_sarama.WithLoopBatchFunction(
// 			func(c context.Context, m []message.Message[message.Bytes, message.Bytes]) ([]message.Message[message.Bytes, message.Bytes], error) {
// 				executionCount += 1
// 				assert.Equal(1, len(m))
// 				return nil, errors.New("mocked error 2")
// 			},
// 		),
// 	)

// 	// mock setup

// 	session.EXPECT().MarkMessage(gomock.Eq(&sarama.ConsumerMessage{Partition: 1}), gomock.Any()).Times(0)

// 	// execute test

// 	go func() {
// 		err := consumerBatchLoop.ConsumeClaim(session, ConsumerGroupClaimForTest{messages: messages})
// 		assert.NotNil(err)
// 		assert.Equal("mocked error 2", err.Error())
// 		if err == nil {
// 			completed <- true
// 		} else {
// 			completed <- false
// 		}
// 	}()

// 	time.Sleep(time.Millisecond)
// 	assert.Equal(0, len(completed))
// 	messages <- &sarama.ConsumerMessage{Partition: 1}
// 	time.Sleep(101 * time.Millisecond)
// 	close(messages)
// 	time.Sleep(101 * time.Millisecond)
// 	assert.Equal(1, len(completed))

// 	result := <-completed
// 	assert.False(result)
// 	assert.Equal(1, executionCount)
// }
