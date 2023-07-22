package runtime_sarama

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"go.uber.org/zap"
)

// constructor
func NewConsumer(configurations ...runtime.Configuration[*Consumer]) runtime.Runtime {
	consumer := &Consumer{
		SaramaConfiguration: DefaultConfiguration(),
	}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return runtime.NewLoop(consumer)
}

type Consumer struct {
	// required
	Topics              []string
	Brokers             []string
	SaramaConfiguration *sarama.Config
	Handler             ConsumerLoop

	// not required
	GroupName string

	// set in start
	Group sarama.ConsumerGroup
}

func (r *Consumer) Start() error {
	// basic validations
	if r == nil {
		return ErrConsumerIsNil
	}
	if r.Handler == nil {
		return ErrConsumerLoopIsNil
	}
	if r.SaramaConfiguration == nil {
		return ErrConsumerSaramaConfigurationIsNil
	}
	if len(r.Topics) == 0 {
		return ErrConsumerTopicsEmpty
	}
	if len(r.Brokers) == 0 {
		return ErrConsumerBrokersEmpty
	}

	logger.Info("starting sarama consumer")

	// set consumer group default
	if len(r.GroupName) == 0 {
		r.GroupName = "flows-" + uuid.New().String()
	}
	logger.Info("using consumer group name", zap.String("name", r.GroupName))

	// create consumer group
	var groupInitErr error
	r.Group, groupInitErr = sarama.NewConsumerGroup(r.Brokers, r.GroupName, r.SaramaConfiguration)
	if groupInitErr != nil {
		return groupInitErr
	}

	logger.Info("started sarama consumer")

	return nil
}

func (r *Consumer) Stop() {
	logger.Info("stopping sarama consumer")
	if r.Group != nil {
		r.Group.Close()
	}
	logger.Info("stopped sarama consumer")
}

func (r *Consumer) Loop(ctx context.Context, cancel context.CancelFunc) error {
	// `Consume` should be called inside an infinite loop, when a
	// server-side rebalance happens, the consumerGroup session will need to be
	// recreated to get the new claims
	if len(r.Group.Errors()) > 0 {
		return <-r.Group.Errors()
	}

	if err := r.Group.Consume(ctx, r.Topics, r.Handler); err != nil {
		if err.Error() == "kafka: tried to use a consumer group that was closed" {
			cancel()
			return nil
		}

		logger.ErrorErr("consume group run error", err)
		return err
	}
	return nil
}

type ConsumerLoop interface {
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	Start() error
	Stop()
}

// Errors
var (
	ErrConsumerIsNil                    = errors.New("consumer is nil")
	ErrConsumerLoopIsNil                = errors.New("consumer loop is nil")
	ErrConsumerSaramaConfigurationIsNil = errors.New("consumer sarama configuration is nil")
	ErrConsumerTopicsEmpty              = errors.New("consumer topics are empty")
	ErrConsumerBrokersEmpty             = errors.New("consumer brokers are empty")
)
