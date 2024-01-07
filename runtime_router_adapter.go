package flows

import (
	"context"

	"github.com/hjwalt/flows/adapter"
	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/runtime_sarama"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

// Wiring configuration
type RouterAdapterConfiguration[Request any, InputKey any, InputValue any] struct {
	Name                       string
	Path                       string
	ProduceTopic               flow.Topic[InputKey, InputValue]
	ProduceBroker              string
	RequestBodyFormat          format.Format[Request]
	RequestMapFunction         stateless.OneToOneFunction[structure.Bytes, Request, InputKey, InputValue]
	HttpPort                   int
	KafkaProducerConfiguration []runtime.Configuration[*runtime_sarama.Producer]
	RouteConfiguration         []runtime.Configuration[*runtime_chi.Runtime[context.Context]]
}

func (c RouterAdapterConfiguration[Request, InputKey, InputValue]) Register(ci inverse.Container) {
	RegisterRouteConfig(
		ci,
		c.RouteConfiguration...,
	)
	RegisterProducerRoute(
		ci,
		"POST",
		c.Path,
		adapter.RouteProduceTopicBodyMapConvert(
			c.RequestMapFunction,
			c.RequestBodyFormat,
			c.ProduceTopic,
		),
	)

	// RUNTIME

	RegisterProducer(
		ci,
		c.ProduceBroker,
		c.KafkaProducerConfiguration,
	)
	RegisterRoute(
		ci,
		c.HttpPort,
		c.RouteConfiguration,
	)
}
