package flows

import (
	"context"

	"github.com/hjwalt/flows/adapter"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	QualifierRoute = "QualifierRoute"
)

func RegisterRoute(
	container inverse.Container,
	port int,
	configs []runtime.Configuration[*runtime_chi.Runtime[context.Context]],
) {

	resolver := runtime.NewResolver[*runtime_chi.Runtime[context.Context], runtime.Runtime](
		QualifierRoute,
		container,
		true,
		runtime_chi.New[context.Context],
	)

	if port == 0 {
		port = 8081
	}

	resolver.AddConfigVal(runtime_chi.WithPort[context.Context](port))
	resolver.AddConfigVal(runtime_chi.WithHttpHandler[context.Context]("/prometheus", "GET", promhttp.Handler()))

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRoute, container)
}

// ===================================

func RegisterRouteConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_chi.Runtime[context.Context]]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierRoute), config)
	}
}

func RegisterProducerRoute(ci inverse.Container, method string, path string, bodyMap stateless.OneToOneFunction[structure.Bytes, structure.Bytes, structure.Bytes, structure.Bytes]) {
	inverse.GenericAdd(
		ci,
		runtime.QualifierConfig(QualifierRoute),
		func(ctx context.Context, c inverse.Container) (runtime.Configuration[*runtime_chi.Runtime[context.Context]], error) {
			producer, getProducerError := GetKafkaProducer(ctx, ci)
			if getProducerError != nil {
				return nil, getProducerError
			}

			handlerFunction := adapter.NewRouteProducer(
				adapter.WithRouteProducerRuntime(producer),
				adapter.WithRouteBodyMap(bodyMap),
			)

			return runtime_chi.WithHttpHandler[context.Context](path, method, handlerFunction), nil
		},
	)
}
