package flows

import (
	"context"

	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_neo4j"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
// Hard to describe without a execution tree graph, but the rough idea is as follows
// consumer -> retry
// retry -> producer
// producer -> stateless switch
// stateless switch -> source to intermediate (for each source topic)
// stateless switch -> intermediate to join
// intermediate to join -> transaction
// transaction -> offset deduplication
// offset deduplication -> stateful switch
// stateful switch -> stateful function(s) (for each source topic)

type JoinPostgresqlFunctionConfiguration struct {
	Name                       string
	StatefulFunctions          map[string]stateful.SingleFunction
	PersistenceIdFunctions     map[string]stateful.PersistenceIdFunction[[]byte, []byte]
	IntermediateTopicName      string
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	PostgresConnectionString   string
	PostgresTable              string
	PostgresqlConfiguration    []runtime.Configuration[*runtime_bun.PostgresqlConnection]
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
	Neo4jConfiguration         []runtime.Configuration[*runtime_neo4j.Neo4JConnectionBasicAuth]
}

func (c JoinPostgresqlFunctionConfiguration) Register(ci inverse.Container) {
	statefulTopicSwitchConfigurations := []runtime.Configuration[*stateful.TopicSwitch]{}
	persistenceIdConfigurations := []runtime.Configuration[*stateful.PersistenceIdSwitch]{}

	// generating stateful switches for persistence id and stateful function
	for topic, statefulFn := range c.StatefulFunctions {
		persistenceIdFn, persistenceIdFnExists := c.PersistenceIdFunctions[topic]
		if !persistenceIdFnExists {
			// TODO: not so silent failure
			continue
		}

		sourceToIntermediateMap := join.NewSourceToIntermediateMap(
			join.WithSourceToIntermediateMapIntermediateTopic(c.IntermediateTopicName),
			join.WithSourceToIntermediateMapPersistenceIdFunction(persistenceIdFn),
		)

		RegisterStatelessSingleFunctionWithKey(
			ci,
			topic,
			sourceToIntermediateMap,
			persistenceIdFn,
		)

		statefulTopicSwitchConfigurations = append(statefulTopicSwitchConfigurations, stateful.WithTopicSwitchFunction(topic, statefulFn))
		persistenceIdConfigurations = append(persistenceIdConfigurations, stateful.WithPersistenceIdSwitchPersistenceIdFunction(topic, persistenceIdFn))
	}

	RegisterJoinStatefulFunction(
		ci,
		c.IntermediateTopicName,
		c.PostgresTable,
		stateful.NewTopicSwitch(statefulTopicSwitchConfigurations...),
		stateful.NewPersistenceIdSwitch(persistenceIdConfigurations...),
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
}
