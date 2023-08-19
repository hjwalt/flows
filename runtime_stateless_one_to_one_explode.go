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

type StatelessOneToOneExplodeConfiguration[IK any, IV any, OK any, OV any] struct {
	Name                       string
	InputTopic                 topic.Topic[IK, IV]
	OutputTopic                topic.Topic[OK, OV]
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
			router.WithFlowStatelessOneToOne(c.InputTopic, c.OutputTopic),
		),
	}
	if len(c.RouteConfiguration) > 0 {
		routeConfigs = append(routeConfigs, c.RouteConfiguration...)
	}

	statelessFunctionConfiguration := StatelessSingleFunctionConfiguration{
		StatelessFunction:          stateless.ConvertTopicOneToOneExplode(c.Function, c.InputTopic, c.OutputTopic),
		KafkaProducerConfiguration: kafkaProducerConfigs,
		KafkaConsumerConfiguration: kafkaConsumerConfigs,
		RouteConfiguration:         routeConfigs,
		RetryConfiguration:         c.RetryConfiguration,
	}

	statelessFunctionConfiguration.Register()
}

func (c StatelessOneToOneExplodeConfiguration[IK, IV, OK, OV]) Runtime() runtime.Runtime {
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(),
	}
}
