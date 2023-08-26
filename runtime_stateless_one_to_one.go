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

type StatelessOneToOneConfiguration[IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 topic.Topic[IK, IV]
	OutputTopic                topic.Topic[OK, OV]
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
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
	)

	statelessFunctionConfiguration := StatelessSingleFunctionConfiguration{
		StatelessFunction:          stateless.ConvertTopicOneToOne(c.Function, c.InputTopic, c.OutputTopic),
		KafkaProducerConfiguration: c.KafkaProducerConfiguration,
		KafkaConsumerConfiguration: c.KafkaConsumerConfiguration,
		RouteConfiguration:         c.RouteConfiguration,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statelessFunctionConfiguration.Register()
}

func (c StatelessOneToOneConfiguration[IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
