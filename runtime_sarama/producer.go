package runtime_sarama

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
)

func NewProducer(configurations ...runtime.Configuration[*Producer]) runtime.Producer {
	producer := &Producer{}
	for _, configuration := range configurations {
		producer = configuration(producer)
	}
	return producer
}

type Producer struct {
	// required
	Brokers             []string
	SaramaConfiguration *sarama.Config
	Controller          runtime.Controller

	// set in start
	Producer sarama.SyncProducer
}

func (p *Producer) Produce(c context.Context, sources []message.Message[message.Bytes, message.Bytes]) error {
	if len(sources) == 0 {
		return nil
	}
	mappedMessages, err := ToProducerMessages(sources)
	if err != nil {
		return err
	}
	if len(mappedMessages) > 0 {
		return p.Producer.SendMessages(mappedMessages)
	} else {
		return nil
	}
}

func (p *Producer) Start() error {

	// basic validations
	if p == nil {
		return errors.New("producer is nil")
	}
	if p.Controller == nil {
		return errors.New("producer controller is nil")
	}
	if p.SaramaConfiguration == nil {
		return errors.New("producer sarama configuration is nil")
	}
	if len(p.Brokers) == 0 {
		return errors.New("producer brokers are empty")
	}

	logger.Info("starting sarama producer")

	// create producer
	var producerCreateErr error
	p.Producer, producerCreateErr = sarama.NewSyncProducer(p.Brokers, p.SaramaConfiguration)
	if producerCreateErr != nil {
		return producerCreateErr
	}

	logger.Info("started sarama producer")

	// mark started in controller
	p.Controller.Started()

	return nil
}

func (p *Producer) Stop() {
	logger.Info("stopping sarama producer")

	if p.Producer != nil {
		p.Producer.Close()
	}

	p.Controller.Stopped()
	logger.Info("stopped sarama producer")
}
