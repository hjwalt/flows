package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/materialise"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
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
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
	Neo4jConfiguration         []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]
}

func (c MaterialisePostgresqlOneToOneFunctionConfiguration[S, IK, IV]) Register(ci inverse.Container) {
	RegisterMaterialiseFunction(
		ci,
		c.InputTopic.Name(),
		materialise.ConvertOneToOne(c.Function, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
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
			runtime_neo4j.InsertTable(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.Name)
			runtime_neo4j.InsertRelationMaterialise(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.Name, c.InputTopic, c.Name)
		}
	}
}
