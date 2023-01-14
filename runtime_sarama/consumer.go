package runtime_sarama

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

// constructor
func NewConsumer(configurations ...runtime.Configuration[*Consumer]) *Consumer {
	consumer := &Consumer{}
	for _, configuration := range configurations {
		consumer = configuration(consumer)
	}
	return consumer
}

// implementation
type Consumer struct {
	// required
	Topics              []string
	Brokers             []string
	SaramaConfiguration *sarama.Config
	Loop                ConsumerLoop
	Controller          runtime.Controller

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
		return errors.New("consumer is nil")
	}
	if c.Controller == nil {
		return errors.New("consumer controller is nil")
	}
	if c.Loop == nil {
		return errors.New("consumer loop is nil")
	}
	if c.SaramaConfiguration == nil {
		return errors.New("consumer sarama configuration is nil")
	}
	if len(c.Topics) == 0 {
		return errors.New("consumer topics are empty")
	}
	if len(c.Brokers) == 0 {
		return errors.New("consumer brokers are empty")
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

	// mark started in controller
	c.Controller.Started()

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

func (c *Consumer) Run() {
	logger.Info("sarama consumer run start")
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumerGroup session will need to be
		// recreated to get the new claims
		if err := c.Group.Consume(c.Context, c.Topics, c); err != nil {
			logger.ErrorErr("consume group run error", err)
			c.Cancel()
			break
		}
		// check if context was cancelled, signaling that the consumerGroup should stop
		if c.Context.Err() != nil {
			break
		}
	}

	// mark stopped when exiting loop
	c.Controller.Stopped()
	logger.Info("sarama consumer run end")
}

// Setup is run at the beginning of a new session, before ConsumeClaim, marks the consumer group as ready
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumerGroup loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	claimErr := c.Loop.Loop(session, claim)
	if claimErr != nil {
		logger.Error("consume claim error", zap.Error(claimErr))
	}
	return claimErr
}

type ConsumerLoop interface {
	Loop(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
	Start() error
	Stop()
}
