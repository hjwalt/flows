package flows

import (
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_retry"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/runtime"
)

type StatelessOneToTwoConfiguration[IK any, IV any, OK1 any, OV1 any, OK2 any, OV2 any] struct {
	Name                       string
	InputTopic                 topic.Topic[IK, IV]
	OutputTopicOne             topic.Topic[OK1, OV1]
	OutputTopicTwo             topic.Topic[OK2, OV2]
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
	RegisterConsumerConfig(
		runtime_sarama.WithConsumerBroker(c.InputBroker),
		runtime_sarama.WithConsumerTopic(c.InputTopic.Topic()),
		runtime_sarama.WithConsumerGroupName(c.Name),
	)
	RegisterProducerConfig(
		runtime_sarama.WithProducerBroker(c.OutputBroker),
	)
	RegisterRouteConfig(
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToTwo(c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		),
	)

	statelessFunctionConfiguration := StatelessSingleFunctionConfiguration{
		StatelessFunction:          stateless.ConvertTopicOneToTwo(c.Function, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		KafkaProducerConfiguration: c.KafkaProducerConfiguration,
		KafkaConsumerConfiguration: c.KafkaConsumerConfiguration,
		RouteConfiguration:         c.RouteConfiguration,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statelessFunctionConfiguration.Register()
}

func (c StatelessOneToTwoConfiguration[IK, IV, OK1, OV1, OK2, OV2]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
