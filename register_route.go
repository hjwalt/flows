package flows

import (
	"context"

	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRouteConfiguration = "QualifierRouteConfiguration"
	QualifierRoute              = "QualifierRoute"
)

func RegisterRoute(
	port int,
	configs []runtime.Configuration[*runtime_bunrouter.Router],
) {
	inverse.RegisterConfiguration(QualifierRouteConfiguration, runtime_bunrouter.WithRouterPort(port))
	inverse.RegisterConfiguration(QualifierRouteConfiguration, runtime_bunrouter.WithRouterPrometheus())
	inverse.RegisterInstances(QualifierRouteConfiguration, configs)
	inverse.Register(QualifierRouteConfiguration, ResolveRouteProducer)
	inverse.RegisterWithConfigurationRequired[*runtime_bunrouter.Router](QualifierRoute, QualifierRouteConfiguration, runtime_bunrouter.NewRouter)

	RegisterRuntime(QualifierRoute)
}

func RegisterRouteConfig(config ...runtime.Configuration[*runtime_bunrouter.Router]) {
	inverse.RegisterInstances(QualifierRouteConfiguration, config)
}

func ResolveRouteProducer(ctx context.Context) (runtime.Configuration[*runtime_bunrouter.Router], error) {
	producer, getProducerError := GetKafkaProducer(ctx)
	if getProducerError != nil {
		return nil, getProducerError
	}
	return runtime_bunrouter.WithRouterProducer(producer), nil
}
