package flows

import (
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime_bunrouter"
	"github.com/hjwalt/runway/runtime"
)

func RouteRuntime(producer message.Producer, routeConfiguration []runtime.Configuration[*runtime_bunrouter.Router]) runtime.Runtime {

	routeConfig := append(
		make([]runtime.Configuration[*runtime_bunrouter.Router], 0),
		runtime_bunrouter.WithRouterPrometheus(),
		runtime_bunrouter.WithRouterProducer(producer),
	)
	routeConfig = append(
		routeConfig,
		routeConfiguration...,
	)
	return runtime_bunrouter.NewRouter(routeConfig...)
}
