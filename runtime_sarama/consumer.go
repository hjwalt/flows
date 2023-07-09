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
	consumer := &Consumer{}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return runtime.NewRunner(consumer)
}

// implementation
type Consumer struct {
	// required
	Topics              []string
	Brokers             []string
	SaramaConfiguration *sarama.Config
	Loop                ConsumerLoop
	// Controller          runtime.Controller

	// not required
	GroupName string

	// set in start
	Group   sarama.ConsumerGroup
	Context context.Context
	Cancel  context.CancelFunc
}

func (c *Consumer) Start() error {

	// basic validations
	if c == nil {
		return ErrConsumerIsNil
	}
	if c.Loop == nil {
		return ErrConsumerLoopIsNil
	}
	if c.SaramaConfiguration == nil {
		return ErrConsumerSaramaConfigurationIsNil
	}
	if len(c.Topics) == 0 {
		return ErrConsumerTopicsEmpty
	}
	if len(c.Brokers) == 0 {
		return ErrConsumerBrokersEmpty
	}

	logger.Info("starting sarama consumer")

	// set consumer group default
	if len(c.GroupName) == 0 {
		c.GroupName = "flows-" + uuid.New().String()
	}
	logger.Info("using consumer group name", zap.String("name", c.GroupName))

	// create start and stop context
	c.Context, c.Cancel = context.WithCancel(context.Background())

	// create consumer group
	var groupInitErr error
	c.Group, groupInitErr = sarama.NewConsumerGroup(c.Brokers, c.GroupName, c.SaramaConfiguration)
	if groupInitErr != nil {
		return groupInitErr
	}

	// start loop
	go c.Run()
	logger.Info("started sarama consumer")

	return nil
}

func (c *Consumer) Stop() {
	logger.Info("stopping sarama consumer")

	c.Cancel()

	if c.Group != nil {
		// stopping consumer group
		c.Group.Close()
	}

	logger.Info("stopped sarama consumer")
}

func (c *Consumer) Run() error {
	logger.Info("sarama consumer run start")
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumerGroup session will need to be
		// recreated to get the new claims
		if err := c.Group.Consume(c.Context, c.Topics, c.Loop); err != nil {
			logger.ErrorErr("consume group run error", err)
			c.Cancel()
			break
		}
		// check if context was cancelled, signaling that the consumerGroup should stop
		if c.Context.Err() != nil {
			break
		}
	}

	logger.Info("sarama consumer run end")
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
