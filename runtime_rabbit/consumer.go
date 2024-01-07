package runtime_rabbit

import (
	"context"
	"errors"
	"time"

	"github.com/hjwalt/flows/task"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"github.com/rabbitmq/amqp091-go"
)

// constructor
func NewConsumer(configurations ...runtime.Configuration[*Consumer]) runtime.Runtime {
	loop := &Consumer{
		Name:         "tasks",
		QueueDurable: true,
	}
	for _, configuration := range configurations {
		loop = configuration(loop)
	}
	return runtime.NewLoop(loop)
}

// configuration
func WithConsumerName(name string) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Name = name
		return c
	}
}

func WithConsumerConnectionString(connectionString string) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.ConnectionString = connectionString
		return c
	}
}

func WithConsumerQueueName(queueName string) runtime.Configuration[*Consumer] {
	return func(p *Consumer) *Consumer {
		p.QueueName = queueName
		return p
	}
}

func WithConsumerQueueDurable(durable bool) runtime.Configuration[*Consumer] {
	return func(p *Consumer) *Consumer {
		p.QueueDurable = durable
		return p
	}
}

func WithConsumerHandler(handler task.Executor[structure.Bytes]) runtime.Configuration[*Consumer] {
	return func(c *Consumer) *Consumer {
		c.Handler = handler
		return c
	}
}

// implementation
type Consumer struct {
	Name             string
	ConnectionString string // amqp://guest:guest@localhost:5672/
	QueueName        string
	QueueDurable     bool
	Handler          task.Executor[structure.Bytes]

	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      *amqp091.Queue
	messages   <-chan amqp091.Delivery
}

func (r *Consumer) Start() error {
	config := amqp091.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Properties: amqp091.Table{
			"connection_name": r.Name,
		},
	}

	logger.Debug("rabbit consumer starting")

	if conn, err := amqp091.DialConfig(r.ConnectionString, config); err != nil {
		return errors.Join(err, ErrRabbitConnection)
	} else {
		r.connection = conn
	}

	if ch, err := r.connection.Channel(); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		r.channel = ch
	}

	if err := r.channel.Confirm(false); err != nil {
		return errors.Join(err, ErrRabbitConfirmMode)
	}

	if err := r.channel.Qos(1, 0, false); err != nil {
		return errors.Join(err, ErrRabbitPrefetch)
	}

	if q, err := r.channel.QueueDeclare(r.QueueName, r.QueueDurable, false, false, false, nil); err != nil {
		return errors.Join(err, ErrRabbitQueue)
	} else {
		r.queue = &q
	}

	if msgs, err := r.channel.Consume(
		r.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		r.messages = msgs
	}

	logger.Debug("rabbit consumer started")

	return nil
}

func (r *Consumer) Stop() {
	logger.Debug("rabbit consumer stopping")
	r.channel.Close()
	r.connection.Close()
	logger.Debug("rabbit consumer stopped")
}

func (r *Consumer) Loop(ctx context.Context, cancel context.CancelFunc) error {

	m, hasMore := <-r.messages
	if !hasMore {
		cancel()
		return nil
	}

	t := task.Message[structure.Bytes]{
		Channel:   r.queue.Name,
		Value:     m.Body,
		Headers:   m.Headers,
		Timestamp: m.Timestamp,
	}

	if err := r.Handler(ctx, t); err != nil {
		return errors.Join(err, ErrRabbitConsume)
	}

	m.Ack(false)

	return nil
}

var (
	ErrRabbitPrefetch = errors.New("rabbit prefetch setting")
	ErrRabbitMessages = errors.New("rabbit messages start consume")
	ErrRabbitConsume  = errors.New("rabbit messages consume")
)
