package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

type StatelessOneToOneConfiguration[IK any, IV any, OK any, OV any] struct {
	Container                  inverse.Container
	Name                       string
	InputTopic                 flow.Topic[IK, IV]
	OutputTopic                flow.Topic[OK, OV]
	Function                   stateless.OneToOneFunction[IK, IV, OK, OV]
	InputBroker                string
	OutputBroker               string
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	KafkaConsumerConfiguration []runtime.Configuration[*runtime_sarama.Consumer]
	RetryConfiguration         []runtime.Configuration[*runtime_retry.Retry]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) Register() {
	RegisterStatelessSingleFunction(
		c.Container,
		c.InputTopic.Name(),
		stateless.ConvertTopicOneToOne(c.Function, c.InputTopic, c.OutputTopic),
	)
	RegisterRouteConfig(
		c.Container,
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
	)
}

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) RegisterRuntime() {
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

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.RegisterRuntime()
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(
			c.Container,
		),
	}
}

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) Inverse() inverse.Container {
	return c.Container
}
