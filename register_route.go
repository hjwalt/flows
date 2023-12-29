package flows

import (
	"context"

	"github.com/hjwalt/flows/adapter"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
)

const (
	QualifierRoute = "QualifierRoute"
)

func RegisterRoute(
	container inverse.Container,
	port int,
	configs []runtime.Configuration[*runtime_bunrouter.Router],
) {

	resolver := runtime.NewResolver[*runtime_bunrouter.Router, runtime.Runtime](
		QualifierRoute,
		container,
		true,
		runtime_bunrouter.NewRouter,
	)

	if port == 0 {
		port = 8081
	}

	resolver.AddConfigVal(runtime_bunrouter.WithRouterPort(port))
	resolver.AddConfigVal(runtime_bunrouter.WithRouterPrometheus())

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRoute, container)
}

// ===================================

func RegisterRouteConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_bunrouter.Router]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierRoute), config)
	}
}

func RegisterProducerRoute(ci inverse.Container, method string, path string, bodyMap stateless.OneToOneFunction[structure.Bytes, structure.Bytes, structure.Bytes, structure.Bytes]) {
	inverse.GenericAdd(
		ci,
		runtime.QualifierConfig(QualifierRoute),
		func(ctx context.Context, c inverse.Container) (runtime.Configuration[*runtime_bunrouter.Router], error) {
			producer, getProducerError := GetKafkaProducer(ctx, ci)
			if getProducerError != nil {
				return nil, getProducerError
			}

			handlerFunction := adapter.NewRouteProducer(
				adapter.WithRouteProducerRuntime(producer),
				adapter.WithRouteBodyMap(bodyMap),
			)

			return runtime_bunrouter.WithRouterHttpHandler(method, path, handlerFunction.ServeHTTP), nil
		},
	)
}
