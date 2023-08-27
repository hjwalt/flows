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
type StatefulPostgresqlOneToOneFunctionConfiguration[S any, IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopic                flow.Topic[OK, OV]
	Function                   stateful.OneToOneFunction[S, IK, IV, OK, OV]
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

func (c StatefulPostgresqlOneToOneFunctionConfiguration[S, IK, IV, OK, OV]) Register() {
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
			router.WithFlowStatefulOneToOne(c.InputTopic, c.OutputTopic, c.PostgresTable),
		),
	)

	statefulFunctionConfiguration := StatefulPostgresqlFunctionConfiguration{
		PersistenceTableName:       c.PostgresTable,
		PersistenceIdFunction:      stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		StatefulFunction:           stateful.ConvertTopicOneToOne(c.Function, c.StateFormat, c.InputTopic, c.OutputTopic),
		PostgresqlConfiguration:    c.PostgresqlConfiguration,
		KafkaProducerConfiguration: c.KafkaProducerConfiguration,
		KafkaConsumerConfiguration: c.KafkaConsumerConfiguration,
		RouteConfiguration:         c.RouteConfiguration,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statefulFunctionConfiguration.Register()
}

func (c StatefulPostgresqlOneToOneFunctionConfiguration[S, IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
