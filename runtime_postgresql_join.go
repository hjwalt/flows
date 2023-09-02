package flows

import (
	"github.com/hjwalt/flows/join"
	"github.com/hjwalt/flows/runtime_bun"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
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
	Container                  inverse.Container
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
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c JoinPostgresqlFunctionConfiguration) Register() {
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
			c.Container,
			topic,
			sourceToIntermediateMap,
			persistenceIdFn,
		)

		statefulTopicSwitchConfigurations = append(statefulTopicSwitchConfigurations, stateful.WithTopicSwitchFunction(topic, statefulFn))
		persistenceIdConfigurations = append(persistenceIdConfigurations, stateful.WithPersistenceIdSwitchPersistenceIdFunction(topic, persistenceIdFn))
	}

	RegisterJoinStatefulFunction(
		c.Container,
		c.IntermediateTopicName,
		c.PostgresTable,
		stateful.NewTopicSwitch(statefulTopicSwitchConfigurations...),
		stateful.NewPersistenceIdSwitch(persistenceIdConfigurations...),
	)
}

func (c JoinPostgresqlFunctionConfiguration) RegisterRuntime() {
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

func (c JoinPostgresqlFunctionConfiguration) Runtime() runtime.Runtime {
	c.RegisterRuntime()
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(
			c.Container,
		),
	}
}

func (c JoinPostgresqlFunctionConfiguration) Inverse() inverse.Container {
	return c.Container
}
