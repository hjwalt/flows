package flows

import (
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/flows/topic"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/runtime"
)

// Wiring configuration
type RouterAdapterConfiguration[Request any, InputKey any, InputValue any] struct {
	Name                       string
	ProduceTopic               topic.Topic[InputKey, InputValue]
	ProduceBroker              string
	RequestBodyFormat          format.Format[Request]
	RequestMapFunction         stateless.OneToOneFunction[message.Bytes, Request, InputKey, InputValue]
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
	AdditionalRuntimes         []runtime.Runtime
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Runtime() runtime.Runtime {

	// producer configs
	kafkaProducerConfigs := []runtime.Configuration[*runtime_sarama.Producer]{
		runtime_sarama.WithProducerBroker(c.ProduceBroker),
	}
	if len(c.KafkaProducerConfiguration) > 0 {
		kafkaProducerConfigs = append(kafkaProducerConfigs, c.KafkaProducerConfiguration...)
	}

	// route configs
	routeConfigs := []runtime.Configuration[*runtime_bunrouter.Router]{
		runtime_bunrouter.WithRouterPort(c.HttpPort),
		runtime_bunrouter.WithRouterProducerHandler(
			runtime_bunrouter.POST,
			"/"+c.Name,
			router.RouteProduceTopicBodyMapConvert(
				c.RequestMapFunction,
				c.RequestBodyFormat,
				c.ProduceTopic,
			),
		),
	}
	if len(c.RouteConfiguration) > 0 {
		routeConfigs = append(routeConfigs, c.RouteConfiguration...)
	}

	routerConfiguration := RouterConfiguration{
		KafkaProducerConfiguration: kafkaProducerConfigs,
		RouteConfiguration:         routeConfigs,
		AdditionalRuntimes:         c.AdditionalRuntimes,
	}

	return routerConfiguration.Runtime()
}
