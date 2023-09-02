package flows

import (
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

// Wiring configuration
type RouterAdapterConfiguration[Request any, InputKey any, InputValue any] struct {
	Container                  inverse.Container
	Name                       string
	ProduceTopic               flow.Topic[InputKey, InputValue]
	ProduceBroker              string
	RequestBodyFormat          format.Format[Request]
	RequestMapFunction         stateless.OneToOneFunction[structure.Bytes, Request, InputKey, InputValue]
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RouteConfiguration         []runtime.Configuration[*runtime_bunrouter.Router]
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Register() {
	RegisterRouteConfig(
		c.Container,
		runtime_bunrouter.WithRouterProducerHandler(
			runtime_bunrouter.POST,
			"/"+c.Name,
			router.RouteProduceTopicBodyMapConvert(
				c.RequestMapFunction,
				c.RequestBodyFormat,
				c.ProduceTopic,
			),
		),
	)
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) RegisterRuntime() {
	RegisterProducer(
		c.Container,
		c.ProduceBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterRoute(
		c.Container,
		c.HttpPort,
		c.RouteConfiguration,
	)
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Runtime() runtime.Runtime {
	c.RegisterRuntime()
	c.Register()

	return &RuntimeFacade{
		Runtimes: InjectedRuntimes(
			c.Container,
		),
	}
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Inverse() inverse.Container {
	return c.Container
}
