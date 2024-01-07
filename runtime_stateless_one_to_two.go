package flows

import (
	"context"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
)

type StatelessOneToTwoConfiguration[IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopicOne             flow.Topic[OK1, OV1]
	OutputTopicTwo             flow.Topic[OK2, OV2]
	Function                   stateless.OneToTwoFunction[IK, IV, OK1, OV1, OK2, OV2]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
	Neo4jConfiguration         []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]
}

func (c StatelessOneToTwoConfiguration[IK, IV, OK1, OV1, OK2, OV2]) Register(ci inverse.Container) {
	RegisterStatelessSingleFunction(
		ci,
		c.InputTopic.Name(),
		stateless.ConvertTopicOneToTwo(c.Function, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
	)

	// RUNTIME

	RegisterRetry(
		ci,
		c.RetryConfiguration,
	)
	RegisterProducer(
		ci,
		c.OutputBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterConsumer(
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
			runtime_neo4j.InsertTopic(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.OutputTopicOne)
			runtime_neo4j.InsertTopic(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.OutputTopicTwo)
			runtime_neo4j.InsertRelationStateless(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.Name, c.InputTopic, c.OutputTopicOne)
			runtime_neo4j.InsertRelationStateless(neo4jConnection.(runtime_neo4j.Neo4jConnection), c.Name, c.InputTopic, c.OutputTopicTwo)
		}
	}
}
