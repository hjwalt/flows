package runtime_sarama

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/metric"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewBatchLoop(configurations ...runtime.Configuration[*ConsumerBatchedLoop]) ConsumerLoop {
	loop := &ConsumerBatchedLoop{}
	for _, configuration := range configurations {
		loop = configuration(loop)
	}
	return loop
}

// configuration
func WithLoopBatchMaxBufferred(maxBuffered int64) runtime.Configuration[*ConsumerBatchedLoop] {
	return func(cbl *ConsumerBatchedLoop) *ConsumerBatchedLoop {
		cbl.MaxBufferred = maxBuffered
		return cbl
	}
}

func WithLoopBatchMaxDelay(maxDelay time.Duration) runtime.Configuration[*ConsumerBatchedLoop] {
	return func(cbl *ConsumerBatchedLoop) *ConsumerBatchedLoop {
		cbl.MaxDelay = maxDelay
		return cbl
	}
}

func WithLoopBatchFunction(loopFunction stateless.BatchFunction) runtime.Configuration[*ConsumerBatchedLoop] {
	return func(cbl *ConsumerBatchedLoop) *ConsumerBatchedLoop {
		cbl.F = loopFunction
		return cbl
	}
}

func WithLoopBatchPrometheus() runtime.Configuration[*ConsumerBatchedLoop] {
	return func(cbl *ConsumerBatchedLoop) *ConsumerBatchedLoop {
		cbl.metric = metric.PrometheusConsume()
		return cbl
	}
}

// implementation
type ConsumerBatchedLoop struct {
	// batching configuration
	MaxBufferred int64
	MaxDelay     time.Duration
	F            stateless.BatchFunction

	// batched data
	messageToCommit *sarama.ConsumerMessage
	messages        []*sarama.ConsumerMessage

	// batching mechanic
	countReached chan bool
	timerReached <-chan time.Time

	// channel termination mechanic
	messageClaimed chan *sarama.ConsumerMessage
	context        context.Context
	cancel         context.CancelFunc

	// metrics
	metric metric.Consume
}

func (batchConsume *ConsumerBatchedLoop) Loop(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// create start and stop context
	batchConsume.context, batchConsume.cancel = context.WithCancel(context.Background())
	batchConsume.messageClaimed = make(chan *sarama.ConsumerMessage)
	batchConsume.countReached = make(chan bool)

	batchConsume.reset()
	go batchConsume.claimLoop(session, claim)
	for {
		select {
		case saramaMessage := <-batchConsume.messageClaimed:
			logger.Info("read", zap.String("topic", saramaMessage.Topic), zap.Int32("partition", saramaMessage.Partition), zap.Int64("offset", saramaMessage.Offset))
			batchConsume.messages = append(batchConsume.messages, saramaMessage)
			batchConsume.messageToCommit = saramaMessage
			if len(batchConsume.messages) >= int(batchConsume.MaxBufferred) {
				if err := batchConsume.consumeBatch(session); err != nil {
					batchConsume.cancel()
					return err
				}
			}
		case <-batchConsume.timerReached:
			if err := batchConsume.consumeBatch(session); err != nil {
				batchConsume.cancel()
				return err
			}
		}
		if batchConsume.context.Err() != nil {
			return nil
		}
	}
}

func (batchConsume *ConsumerBatchedLoop) reset() {
	batchConsume.messageToCommit = nil
	batchConsume.messages = make([]*sarama.ConsumerMessage, 0)
	batchConsume.timerReached = time.After(batchConsume.MaxDelay)
}

func (batchConsume *ConsumerBatchedLoop) claimLoop(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) {
	// the claim loop channel will stop when the consumer is stopped, triggering end, triggering consume loop end
	for saramaMessage := range claim.Messages() {
		batchConsume.messageClaimed <- saramaMessage
		if batchConsume.context.Err() != nil {
			return
		}
	}
	batchConsume.cancel()
}

func (batchConsume *ConsumerBatchedLoop) consumeBatch(session sarama.ConsumerGroupSession) error {
	if len(batchConsume.messages) == 0 {
		batchConsume.reset()
		return nil
	}

	logger.Info("batch", zap.Int("length", len(batchConsume.messages)))
	sources, err := FromConsumerMessages(batchConsume.messages)
	if err != nil {
		return err
	}
	if _, err := batchConsume.F(context.Background(), sources); err != nil {
		return err
	}
	session.MarkMessage(batchConsume.messageToCommit, "")
	if batchConsume.metric != nil {
		batchConsume.metric.MessagesProcessedIncrement(batchConsume.messageToCommit.Topic, batchConsume.messageToCommit.Partition, int64(len(batchConsume.messages)))
	}
	logger.Info("commit", zap.String("topic", batchConsume.messageToCommit.Topic), zap.Int32("partition", batchConsume.messageToCommit.Partition), zap.Int64("offset", batchConsume.messageToCommit.Offset))

	batchConsume.reset()
	return nil
}

func (consumerSarama *ConsumerBatchedLoop) Start() error {
	return nil
}
func (consumerSarama *ConsumerBatchedLoop) Stop() {
}
