package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
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
	RetryConfiguration         []runtime.Configuration[*runtime.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
	Neo4jConfiguration         []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]
}

func (c StatefulPostgresqlOneToOneFunctionConfiguration[S, IK, IV, OK, OV]) Register(ci inverse.Container) {
	RegisterStatefulFunction(
		ci,
		c.InputTopic.Name(),
		c.PostgresTable,
		stateful.ConvertTopicOneToOne(c.Function, c.StateFormat, c.InputTopic, c.OutputTopic),
		stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
	)
	// RUNTIME

	RegisterPostgresql(
		ci,
		c.Name,
		c.PostgresConnectionString,
		c.PostgresqlConfiguration,
	)
	RegisterRetry(
		ci,
		c.RetryConfiguration,
	)
	RegisterKafkaProducer(
		ci,
		c.OutputBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterKafkaConsumer(
		ci,
		c.Name,
		c.InputBroker,
		c.KafkaConsumerConfiguration,
	)
	RegisterRoute(
		ci,
		c.HttpPort,
		c.RouteConfiguration,
	)

	if len(c.Neo4jConfiguration) > 0 {
		neo4jConnection := runtime_neo4j.NewBasicAuth(c.Neo4jConfiguration...)
		err := neo4jConnection.Start()
		if err != nil {
			logger.ErrorErr("failed to start neo4j", err)
		} else {
			defer neo4jConnection.Stop()

			runtime_neo4j.InsertConstraint(neo4jConnection.(runtime_neo4j.Neo4jConnection))
			runtime_neo4j.InsertTopic(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.InputTopic)
			runtime_neo4j.InsertTopic(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.OutputTopic)
			runtime_neo4j.InsertRelationStateful(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.Name, c.InputTopic, c.OutputTopic)
		}
	}
}
