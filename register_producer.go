package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

var (
	QualifierKafkaProducerConfiguration = "QualifierKafkaProducerConfiguration"
	QualifierKafkaProducer              = "QualifierKafkaProducer"
)

func RegisterProducer(
	broker string,
	configs []runtime.Configuration[*runtime_sarama.Producer],
) {
	RegisterProducerConfig(runtime_sarama.WithProducerBroker(broker))
	RegisterProducerConfig(configs...)
	inverse.RegisterWithConfigurationRequired[*runtime_sarama.Producer](
		QualifierKafkaProducer,
		QualifierKafkaProducerConfiguration,
		runtime_sarama.NewProducer,
	)
	RegisterRuntime(QualifierKafkaProducer)
}

func RegisterProducerConfig(config ...runtime.Configuration[*runtime_sarama.Producer]) {
	inverse.RegisterInstances(QualifierKafkaProducerConfiguration, config)
}

func GetKafkaProducer(ctx context.Context) (flow.Producer, error) {
	return inverse.GetLast[flow.Producer](ctx, QualifierKafkaProducer)
}
