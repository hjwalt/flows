package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type MaterialisePostgresqlOneToOneFunctionConfiguration[S any, IK any, IV any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	Function                   materialise.MapFunction[IK, IV, S]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	PostgresConnectionString   string
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) Register() {
	RegisterPostgresqlConfig(
		runtime_bun.WithApplicationName(c.Name),
		runtime_bun.WithConnectionString(c.PostgresConnectionString),
	)
	RegisterConsumerConfig(
		runtime_sarama.WithConsumerBroker(c.InputBroker),
		runtime_sarama.WithConsumerTopic(c.InputTopic.Name()),
		runtime_sarama.WithConsumerGroupName(c.Name),
	)
	RegisterProducerConfig(
		runtime_sarama.WithProducerBroker(c.OutputBroker),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowMaterialiseOneToOne(c.InputTopic),
		),
	)

	statefulFunctionConfiguration := MaterialisePostgresqlFunctionConfiguration[S]{
		MaterialiseMapFunction:     materialise.ConvertOneToOne(c.Function, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		PostgresqlConfiguration:    c.PostgresqlConfiguration,
		KafkaProducerConfiguration: c.KafkaProducerConfiguration,
		KafkaConsumerConfiguration: c.KafkaConsumerConfiguration,
		RouteConfiguration:         c.RouteConfiguration,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statefulFunctionConfiguration.Register()
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
