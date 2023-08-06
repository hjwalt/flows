package flows

import (
	"context"
	"errors"

	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/runway/inverse"
	"github.com/hjwalt/runway/runtime"
)

const (
	QualifierRouteConfiguration = "QualifierRouteConfiguration"
	QualifierRoute              = "QualifierRoute"
)

func RegisterRoute(config []runtime.Configuration[*runtime_bunrouter.Router]) {
	inverse.RegisterInstances(QualifierRouteConfiguration, config)
	inverse.Register(QualifierRoute, InjectorRoute)
	inverse.Register(QualifierRuntime, InjectorRuntime(QualifierRoute))
}

func InjectorRoute(ctx context.Context) (runtime.Runtime, error) {
	producer, getProducerError := GetKafkaProducer(ctx)
	if getProducerError != nil {
		return nil, getProducerError
	}

	configurations, getConfigurationError := inverse.GetAll[runtime.Configuration[*runtime_bunrouter.Router]](ctx, QualifierRouteConfiguration)
	if getConfigurationError != nil && !errors.Is(getConfigurationError, inverse.ErrNotInjected) {
		return nil, getConfigurationError
	}

	routeConfig := append(
		make([]runtime.Configuration[*runtime_bunrouter.Router], 0),
		runtime_bunrouter.WithRouterPrometheus(),
		runtime_bunrouter.WithRouterProducer(producer),
	)
	routeConfig = append(
		routeConfig,
		configurations...,
	)
	return runtime_bunrouter.NewRouter(routeConfig...), nil
}
