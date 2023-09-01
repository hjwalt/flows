package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatefulPostgresqlOneToTwoFunctionConfiguration[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopicOne             flow.Topic[OK1, OV1]
	OutputTopicTwo             flow.Topic[OK2, OV2]
	Function                   stateful.OneToTwoFunction[S, IK, IV, OK1, OV1, OK2, OV2]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	StateFormat                format.Format[S]
	StateKeyFunction           stateful.PersistenceIdFunction[IK, IV]
	PostgresTable              string
	PostgresConnectionString   string
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatefulPostgresqlOneToTwoFunctionConfiguration[S, IK, IV, OK1, OV1, OK2, OV2]) Register() {
	RegisterPostgresqlStateful(
		c.Name,
		c.PostgresConnectionString,
		c.PostgresTable,
		c.PostgresqlConfiguration,
	)
	RegisterRetry(
		c.RetryConfiguration,
	)
	RegisterProducer2(
		c.OutputBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterConsumer2(
		c.Name,
		c.InputBroker,
		c.KafkaConsumerConfiguration,
	)
	RegisterStatefulFunction(
		c.InputTopic.Name(),
		stateful.ConvertTopicOneToTwo(c.Function, c.StateFormat, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatefulOneToTwo(c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo, c.PostgresTable),
		),
	)
	RegisterRoute2(
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c StatefulPostgresqlOneToTwoFunctionConfiguration[S, IK, IV, OK1, OV1, OK2, OV2]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
