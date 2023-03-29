package runtime_sarama

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewSingleLoop(configurations ...runtime.Configuration[*ConsumerSingleLoop]) ConsumerLoop {
	loop := &ConsumerSingleLoop{}
	for _, configuration := range configurations {
		loop = configuration(loop)
	}
	return loop
}

// configuration
func WithLoopSingleFunction(loopFunction stateless.SingleFunction) runtime.Configuration[*ConsumerSingleLoop] {
	return func(csl *ConsumerSingleLoop) *ConsumerSingleLoop {
		csl.F = loopFunction
		return csl
	}
}

// implementation
type ConsumerSingleLoop struct {
	F stateless.SingleFunction
}

func (consumerSarama *ConsumerSingleLoop) Loop(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for saramaMessage := range claim.Messages() {
		logger.Info("read", zap.String("topic", saramaMessage.Topic), zap.Int32("partition", saramaMessage.Partition), zap.Int64("offset", saramaMessage.Offset))
		source, err := FromConsumerMessage(saramaMessage)
		if err != nil {
			return err
		}
		if _, err := consumerSarama.F(context.Background(), source); err != nil {
			return err
		}
		session.MarkMessage(saramaMessage, "")
		// add to prometheus
		messageProcessedCounter.WithLabelValues(saramaMessage.Topic, fmt.Sprintf("%d", saramaMessage.Partition)).Inc()
		logger.Info("commit", zap.String("topic", saramaMessage.Topic), zap.Int32("partition", saramaMessage.Partition), zap.Int64("offset", saramaMessage.Offset))
	}
	return nil
}
func (consumerSarama *ConsumerSingleLoop) Start() error {
	return nil
}
func (consumerSarama *ConsumerSingleLoop) Stop() {
}
