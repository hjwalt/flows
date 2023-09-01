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
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c StatelessOneToTwoConfiguration[IK, IV, OK1, OV1, OK2, OV2]) Register() {
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
		stateless.ConvertTopicOneToTwo(c.Function, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToTwo(c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		),
	)
	RegisterRoute2(
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c StatelessOneToTwoConfiguration[IK, IV, OK1, OV1, OK2, OV2]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
