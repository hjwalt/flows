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
	AdditionalRuntimes         []runtime.Runtime
}

func (c StatelessOneToTwoConfiguration[IK, IV, OK1, OV1, OK2, OV2]) Runtime() runtime.Runtime {

	// consumer configs
	kafkaConsumerConfigs := []runtime.Configuration[*runtime_sarama.Consumer]{
		runtime_sarama.WithConsumerBroker(c.InputBroker),
		runtime_sarama.WithConsumerTopic(c.InputTopic.Topic()),
		runtime_sarama.WithConsumerGroupName(c.Name),
	}
	if len(c.KafkaConsumerConfiguration) > 0 {
		kafkaConsumerConfigs = append(kafkaConsumerConfigs, c.KafkaConsumerConfiguration...)
	}

	// producer configs
	kafkaProducerConfigs := []runtime.Configuration[*runtime_sarama.Producer]{
		runtime_sarama.WithProducerBroker(c.OutputBroker),
	}
	if len(c.KafkaProducerConfiguration) > 0 {
		kafkaProducerConfigs = append(kafkaProducerConfigs, c.KafkaProducerConfiguration...)
	}

	// route configs
	routeConfigs := []runtime.Configuration[*runtime_bunrouter.Router]{
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterFlow(
			router.WithFlowStatelessOneToTwo(c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		),
	}
	if len(c.RouteConfiguration) > 0 {
		routeConfigs = append(routeConfigs, c.RouteConfiguration...)
	}

	statelessFunctionConfiguration := StatelessSingleFunctionConfiguration{
		StatelessFunction:          stateless.ConvertTopicOneToTwo(c.Function, c.InputTopic, c.OutputTopicOne, c.OutputTopicTwo),
		KafkaProducerConfiguration: kafkaProducerConfigs,
		KafkaConsumerConfiguration: kafkaConsumerConfigs,
		RouteConfiguration:         routeConfigs,
		RetryConfiguration:         c.RetryConfiguration,
		AdditionalRuntimes:         c.AdditionalRuntimes,
	}

	return statelessFunctionConfiguration.Runtime()
}
