package flows

import (
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type StatefulPostgresqlOneToTwoFunctionConfiguration[S any, IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] struct {
	Name                       string
	InputTopic                 topic.Topic[IK, IV]
	OutputTopicOne             topic.Topic[OK1, OV1]
	OutputTopicTwo             topic.Topic[OK2, OV2]
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
	RegisterPostgresqlConfig(
		runtime_bun.WithApplicationName(c.Name),
		runtime_bun.WithConnectionString(c.PostgresConnectionString),
	)
	RegisterConsumerConfig(
		runtime_sarama.WithConsumerBroker(c.InputBroker),
		runtime_sarama.WithConsumerTopic(c.InputTopic.Topic()),
		runtime_sarama.WithConsumerGroupName(c.Name),
	)
	RegisterProducerConfig(
		runtime_sarama.WithProducerBroker(c.OutputBroker),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatefulOneToTwo(c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo, c.PostgresTable),
		),
	)

	statefulFunctionConfiguration := StatefulPostgresqlFunctionConfiguration{
		PersistenceTableName:       c.PostgresTable,
		PersistenceIdFunction:      stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		StatefulFunction:           stateful.ConvertTopicOneToTwo(c.Function, c.StateFormat, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		PostgresqlConfiguration:    c.PostgresqlConfiguration,
		KafkaProducerConfiguration: c.KafkaProducerConfiguration,
		KafkaConsumerConfiguration: c.KafkaConsumerConfiguration,
		RouteConfiguration:         c.RouteConfiguration,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statefulFunctionConfiguration.Register()
}

func (c StatefulPostgresqlOneToTwoFunctionConfiguration[S, IK, IV, OK1, OV1, OK2, OV2]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
