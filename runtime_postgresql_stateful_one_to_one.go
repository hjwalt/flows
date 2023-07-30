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
type StatefulPostgresqlOneToOneFunctionConfiguration[S any, IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 topic.Topic[IK, IV]
	OutputTopic                topic.Topic[OK, OV]
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

func (c StatefulPostgresqlOneToOneFunctionConfiguration[S, IK, IV, OK, OV]) Runtime() runtime.Runtime {

	// postgres configs
	postgresConfigs := []runtime.Configuration[*runtime_bun.PostgresqlConnection]{
		runtime_bun.WithApplicationName(c.Name),
		runtime_bun.WithConnectionString(c.PostgresConnectionString),
	}
	if len(c.PostgresqlConfiguration) > 0 {
		postgresConfigs = append(postgresConfigs, c.PostgresqlConfiguration...)
	}

	// consumer configs
	kafkaConsumerConfigs := []runtime.Configuration[*runtime_sarama.Consumer]{
		runtime_sarama.WithConsumerBroker(c.InputBroker),
		runtime_sarama.WithConsumerTopic(c.InputTopic.Topic()),
		runtime_sarama.WithConsumerGroupName(c.Name),
	}
	if len(c.KafkaConsumerConfiguration) > 0 {
		kafkaConsumerConfigs = append(kafkaConsumerConfigs, c.KafkaConsumerConfiguration...)
	}

	// producer configs
	kafkaProducerConfigs := []runtime.Configuration[*runtime_sarama.Producer]{
		runtime_sarama.WithProducerBroker(c.OutputBroker),
	}
	if len(c.KafkaProducerConfiguration) > 0 {
		kafkaProducerConfigs = append(kafkaProducerConfigs, c.KafkaProducerConfiguration...)
	}

	// route configs
	routeConfigs := []runtime.Configuration[*runtime_bunrouter.Router]{
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatefulOneToOne(c.InputTopic, c.OutputTopic, c.PostgresTable),
		),
	}
	if len(c.RouteConfiguration) > 0 {
		routeConfigs = append(routeConfigs, c.RouteConfiguration...)
	}

	statefulFunctionConfiguration := StatefulPostgresqlFunctionConfiguration{
		PersistenceTableName:       c.PostgresTable,
		PersistenceIdFunction:      stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		StatefulFunction:           stateful.ConvertTopicOneToOne(c.Function, c.StateFormat, c.InputTopic, c.OutputTopic),
		PostgresqlConfiguration:    postgresConfigs,
		KafkaProducerConfiguration: kafkaProducerConfigs,
		KafkaConsumerConfiguration: kafkaConsumerConfigs,
		RouteConfiguration:         routeConfigs,
		RetryConfiguration:         c.RetryConfiguration,
	}

	return statefulFunctionConfiguration.Runtime()
}