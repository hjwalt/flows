package flows

import (
	"context"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatelessSingleFunctionConfiguration struct {
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	StatelessFunction          stateless.SingleFunction
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessSingleFunctionConfiguration) Runtime() runtime.Runtime {
	RegisterRetry(c.RetryConfiguration)
	RegisterProducerConfig(c.KafkaProducerConfiguration)
	RegisterProducer()
	RegisterConsumerKeyedConfig(c.KafkaConsumerConfiguration)
	RegisterConsumer()
	RegisterRoute(c.RouteConfiguration)
	inverse.RegisterInstance[stateful.PersistenceIdFunction[message.Bytes, message.Bytes]](QualifierKafkaConsumerKeyFunction, statelessBase64PersistenceId)
	inverse.Register[stateless.BatchFunction](QualifierKafkaConsumerBatchFunction, func(ctx context.Context) (stateless.BatchFunction, error) {
		retry, err := GetRetry(ctx)
		if err != nil {
			return nil, err
		}
		producer, err := GetKafkaProducer(ctx)
		if err != nil {
			return nil, err
		}

		wrappedFunction := c.StatelessFunction

		wrappedFunction = stateless.NewSingleRetry(
			stateless.WithSingleRetryNextFunction(wrappedFunction),
			stateless.WithSingleRetryRuntime(retry),
			stateless.WithSingleRetryPrometheus(),
		)

		wrappedBatch := stateless.NewProducerBatchIterateFunction(
			stateless.WithBatchIterateFunctionNextFunction(wrappedFunction),
			stateless.WithBatchIterateFunctionProducer(producer),
			stateless.WithBatchIterateProducerPrometheus(),
		)

		wrappedBatch = stateless.NewBatchRetry(
			stateless.WithBatchRetryNextFunction(wrappedBatch),
			stateless.WithBatchRetryRuntime(retry),
			stateless.WithBatchRetryPrometheus(),
		)

		return wrappedBatch, nil
	})

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}

var (
	base64Format = format.Base64()
	bytesFormat  = format.Bytes()
)

func statelessBase64PersistenceId(ctx context.Context, m message.Message[message.Bytes, message.Bytes]) (string, error) {
	base64Message, conversionErr := message.Convert(
		m,
		bytesFormat,
		bytesFormat,
		base64Format,
		bytesFormat,
	)

	if conversionErr != nil {
		return "", conversionErr
	}

	return base64Message.Key, nil
}
