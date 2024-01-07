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
func NewProducer(configurations ...runtime.Configuration[*Producer]) task.Producer {
	r := &Producer{
		Name: "tasks",
	}
	for _, configuration := range configurations {
		r = configuration(r)
	}
	return r
}

// configuration
func WithProducerName(name string) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.Name = name
		return p
	}
}

func WithProducerConnectionString(connectionString string) runtime.Configuration[*Producer] {
	return func(p *Producer) *Producer {
		p.ConnectionString = connectionString
		return p
	}
}

// implementation
type Producer struct {
	Name             string
	ConnectionString string // amqp://guest:guest@localhost:5672/

	connection *amqp091.Connection
	channel    *amqp091.Channel
}

func (p *Producer) Start() error {
	config := amqp091.Config{
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
		Properties: amqp091.Table{
			"connection_name": p.Name,
		},
	}

	logger.Debug("rabbit producer starting")

	if conn, err := amqp091.DialConfig(p.ConnectionString, config); err != nil {
		return errors.Join(err, ErrRabbitConnection)
	} else {
		p.connection = conn
	}

	if ch, err := p.connection.Channel(); err != nil {
		return errors.Join(err, ErrRabbitChannel)
	} else {
		p.channel = ch
	}

	if err := p.channel.Confirm(false); err != nil {
		return errors.Join(err, ErrRabbitConfirmMode)
	}

	logger.Debug("rabbit producer started")

	return nil
}

func (p *Producer) Stop() {
	logger.Debug("rabbit producer stoppping")

	p.channel.Close()
	p.connection.Close()

	logger.Debug("rabbit producer stopped")
}

func (p *Producer) Produce(c context.Context, t task.Message[structure.Bytes]) error {
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	confirm, err := p.channel.PublishWithDeferredConfirmWithContext(ctx,
		"",        // exchange
		t.Channel, // routing key
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			Headers:      t.Headers,
			ContentType:  "application/octet-stream",
			Body:         t.Value,
		},
	)
	if err != nil {
		return errors.Join(err, ErrRabbitProduce)
	}

	published, err := confirm.WaitContext(ctx)
	if err != nil {
		return errors.Join(err, ErrRabbitProduce)
	}

	if !published {
		return errors.Join(err, ErrRabbitProduce)
	}

	return nil
}

var (
	ErrRabbitConnection          = errors.New("rabbitmq connection failed")
	ErrRabbitChannel             = errors.New("rabbitmq channel failed")
	ErrRabbitQueue               = errors.New("rabbitmq queue declaration failed")
	ErrRabbitProduce             = errors.New("rabbitmq producing")
	ErrRabbitProduceNotConfirmed = errors.New("rabbitmq produce not confirmed")
	ErrRabbitConfirmMode         = errors.New("rabbitmq channel confirm mode failed")
)
