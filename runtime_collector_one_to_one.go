package flows

import (
	"github.com/hjwalt/flows/collect"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type CollectorOneToOneConfiguration[S any, IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopic                flow.Topic[OK, OV]
	Aggregator                 collect.Aggregator[S, IK, IV]
	Collector                  collect.OneToOneCollector[S, OK, OV]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	StateFormat                format.Format[S]
	StateKeyFunction           stateful.PersistenceIdFunction[IK, IV]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
	Neo4jConfiguration         []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]
}

func (c CollectorOneToOneConfiguration[S, IK, IV, OK, OV]) Register(ci inverse.Container) {
	RegisterCollectorFunction(
		ci,
		c.InputTopic.Name(),
		stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		collect.ConvertTopicAggregator(c.Aggregator, c.StateFormat, c.InputTopic),
		collect.ConvertTopicOneToOneCollector(c.Collector, c.StateFormat, c.OutputTopic),
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
}
