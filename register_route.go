package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
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

	resolver.AddConfigVal(runtime_bunrouter.WithRouterPort(port))
	resolver.AddConfigVal(runtime_bunrouter.WithRouterPrometheus())
	resolver.AddConfig(ResolveRouteProducer)

	for _, config := range configs {
		resolver.AddConfigVal(config)
	}

	resolver.Register()

	RegisterRuntime(QualifierRoute, container)
}

func ResolveRouteProducer(ctx context.Context, ci inverse.Container) (runtime.Configuration[*runtime_bunrouter.Router], error) {
	producer, getProducerError := GetKafkaProducer(ctx, ci)
	if getProducerError != nil {
		return nil, getProducerError
	}
	return runtime_bunrouter.WithRouterProducer(producer), nil
}

// ===================================

func RegisterRouteConfig(ci inverse.Container, configs ...runtime.Configuration[*runtime_bunrouter.Router]) {
	for _, config := range configs {
		ci.AddVal(runtime.QualifierConfig(QualifierRoute), config)
	}
}
