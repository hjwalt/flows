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

func RegisterRoute(config []runtime.Configuration[*runtime_bunrouter.Router]) {
	inverse.RegisterConfiguration[*runtime_bunrouter.Router](QualifierRouteConfiguration, runtime_bunrouter.WithRouterPrometheus())
	inverse.Register[runtime.Configuration[*runtime_bunrouter.Router]](QualifierRouteConfiguration, InjectorRouteProducer)
	inverse.RegisterInstances(QualifierRouteConfiguration, config)
	inverse.RegisterWithConfigurationRequired[*runtime_bunrouter.Router](
		QualifierRoute,
		QualifierRouteConfiguration,
		runtime_bunrouter.NewRouter,
	)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierRoute))
}

func InjectorRouteProducer(ctx context.Context) (runtime.Configuration[*runtime_bunrouter.Router], error) {
	producer, getProducerError := GetKafkaProducer(ctx)
	if getProducerError != nil {
		return nil, getProducerError
	}
	return runtime_bunrouter.WithRouterProducer(producer), nil
}
