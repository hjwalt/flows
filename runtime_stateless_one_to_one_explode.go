package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/runtime"
)

type StatelessOneToOneExplodeConfiguration[IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopic                flow.Topic[OK, OV]
	Function                   stateless.OneToOneExplodeFunction[IK, IV, OK, OV]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessOneToOneExplodeConfiguration[IK, IV, OK, OV]) Register() {
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
	RegisterStatelessSingleFunction(
		c.InputTopic.Name(),
		stateless.ConvertTopicOneToOneExplode(c.Function, c.InputTopic, c.OutputTopic),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
	)
	RegisterRoute2(
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c StatelessOneToOneExplodeConfiguration[IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
