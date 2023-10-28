package runtime_sarama

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

func NewProducer(configurations ...runtime.Configuration[*Producer]) flow.Producer {
	producer := &Producer{
		SaramaConfiguration: DefaultConfiguration(),
	}
	for _, configuration := range configurations {
		producer = configuration(producer)
	}
	return producer
}

type Producer struct {
	// required
	Brokers             []string
	SaramaConfiguration *sarama.Config

	// set in start
	Producer sarama.SyncProducer
}

func (p *Producer) Produce(c context.Context, sources []flow.Message[structure.Bytes, structure.Bytes]) error {
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
	if p.SaramaConfiguration == nil {
		return errors.New("producer sarama configuration is nil")
	}
	if len(p.Brokers) == 0 {
		return errors.New("producer brokers are empty")
	}

	logger.Debug("starting sarama producer")

	// create producer
	var producerCreateErr error
	p.Producer, producerCreateErr = sarama.NewSyncProducer(p.Brokers, p.SaramaConfiguration)
	if producerCreateErr != nil {
		return producerCreateErr
	}

	logger.Debug("started sarama producer")

	return nil
}

func (p *Producer) Stop() {
	logger.Debug("stopping sarama producer")

	if p.Producer != nil {
		p.Producer.Close()
	}

	logger.Debug("stopped sarama producer")
}
