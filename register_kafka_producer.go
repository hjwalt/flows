package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

var (
	QualifierKafkaProducer = "QualifierKafkaProducer"
)

func RegisterKafkaProducer(
	container inverse.Container,
	broker string,
	configs []runtime.Configuration[*runtime_sarama.Producer],
) {

	resolver := runtime.NewResolver[*runtime_sarama.Producer, flow.Producer](
		QualifierKafkaProducer,
		container,
		true,
		runtime_sarama.NewProducer,
	)

	resolver.AddConfigVal(runtime_sarama.WithProducerBroker(broker))

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierKafkaProducer, container)
}

func GetKafkaProducer(ctx context.Context, ci inverse.Container) (flow.Producer, error) {
	return inverse.GenericGetLast[flow.Producer](ci, ctx, QualifierKafkaProducer)
}

// ===================================

func RegisterKafkaProducerConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_sarama.Producer]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierKafkaProducer), config)
	}
}
