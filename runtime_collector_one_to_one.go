package flows

import (
	"github.com/hjwalt/flows/collect"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateful"
	"github.com/hjwalt/runway/format"
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
}

func (c CollectorOneToOneConfiguration[S, IK, IV, OK, OV]) Register() {
	RegisterCollectorFunction(
		c.InputTopic.Name(),
		stateful.ConvertPersistenceId(c.StateKeyFunction, c.InputTopic.KeyFormat(), c.InputTopic.ValueFormat()),
		collect.ConvertTopicAggregator(c.Aggregator, c.StateFormat, c.InputTopic),
		collect.ConvertTopicOneToOneCollector(c.Collector, c.StateFormat, c.OutputTopic),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
	)
}

func (c CollectorOneToOneConfiguration[S, IK, IV, OK, OV]) RegisterRuntime() {
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
	RegisterRoute2(
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c CollectorOneToOneConfiguration[S, IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.RegisterRuntime()
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
