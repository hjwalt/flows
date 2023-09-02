package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type MaterialisePostgresqlOneToOneFunctionConfiguration[S any, IK any, IV any] struct {
	Container                  inverse.Container
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
	RegisterMaterialiseFunction(
		c.Container,
		c.InputTopic.Name(),
		materialise.ConvertOneToOne(c.Function, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
	)
	RegisterRouteConfig(
		c.Container,
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowMaterialiseOneToOne(c.InputTopic),
		),
	)
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) RegisterRuntime() {
	RegisterPostgresql(
		c.Container,
		c.Name,
		c.PostgresConnectionString,
		c.PostgresqlConfiguration,
	)
	RegisterRetry(
		c.Container,
		c.RetryConfiguration,
	)
	RegisterProducer(
		c.Container,
		c.OutputBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterConsumer(
		c.Container,
		c.Name,
		c.InputBroker,
		c.KafkaConsumerConfiguration,
	)
	RegisterRoute(
		c.Container,
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) Runtime() runtime.Runtime {
	c.RegisterRuntime()
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(
			c.Container,
		),
	}
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) Inverse() inverse.Container {
	return c.Container
}
