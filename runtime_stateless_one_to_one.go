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

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) Register(ci inverse.Container) {
	RegisterStatelessSingleFunction(
		ci,
		c.InputTopic.Name(),
		stateless.ConvertTopicOneToOne(c.Function, c.InputTopic, c.OutputTopic),
	)
	RegisterRouteConfig(
		ci,
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
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
