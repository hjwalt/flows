package runtime_bunrouter

import (
	"net/http"
	"strings"

	"github.com/hjwalt/flows/flow"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
	"github.com/hjwalt/runway/structure"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uptrace/bunrouter"
	"go.uber.org/zap"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var ports = structure.NewSet[int]()

// constructor
func NewRouter(configurations ...runtime.Configuration[*Router]) runtime.Runtime {
	r := &Router{
		runtimeConfiguration: []runtime.Configuration[*runtime.HttpRunnable]{},
		flowConfiguration:    []runtime.Configuration[*router.RouteFlow]{},
		router:               bunrouter.New(),
	}

	for _, configuration := range configurations {
		r = configuration(r)
	}

	flow := router.NewRouteFlow(
		r.flowConfiguration...,
	)
	r.router.GET("/flow", bunrouter.HTTPHandlerFunc(flow.Handle))

	runtimeConfiguration := append(r.runtimeConfiguration, runtime.HttpWithHandler(r.router))

	return runtime.NewHttp(runtimeConfiguration...)
}

// configuration
func WithRouterPort(port int) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		for ports.Contain(port) {
			port = port + 1
		}

		ports.Add(port)
		r.runtimeConfiguration = append(r.runtimeConfiguration, runtime.HttpWithPort(port))
		return r
	}
}

func WithRouterHttpConfiguration(config runtime.Configuration[*runtime.HttpRunnable]) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.runtimeConfiguration = append(r.runtimeConfiguration, config)
		return r
	}
}

func WithRouterGroup(path string) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.group = r.router.NewGroup(path)
		return r
	}
}

func WithRouterBunHandler(method string, path string, handler bunrouter.HandlerFunc) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		if r.group == nil {
			switch strings.ToUpper(method) {
			case GET:
				r.router.GET(path, handler)
			case POST:
				r.router.POST(path, handler)
			case PUT:
				r.router.PUT(path, handler)
			case DELETE:
				r.router.DELETE(path, handler)
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		} else {
			switch strings.ToUpper(method) {
			case GET:
				r.group.GET(path, handler)
			case POST:
				r.group.POST(path, handler)
			case PUT:
				r.group.PUT(path, handler)
			case DELETE:
				r.group.DELETE(path, handler)
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		}
		return r
	}
}

func WithRouterHttpHandler(method string, path string, handler http.HandlerFunc) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		if r.group == nil {
			switch strings.ToUpper(method) {
			case GET:
				r.router.GET(path, bunrouter.HTTPHandlerFunc(handler))
			case POST:
				r.router.POST(path, bunrouter.HTTPHandlerFunc(handler))
			case PUT:
				r.router.PUT(path, bunrouter.HTTPHandlerFunc(handler))
			case DELETE:
				r.router.DELETE(path, bunrouter.HTTPHandlerFunc(handler))
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		} else {
			switch strings.ToUpper(method) {
			case GET:
				r.group.GET(path, bunrouter.HTTPHandlerFunc(handler))
			case POST:
				r.group.POST(path, bunrouter.HTTPHandlerFunc(handler))
			case PUT:
				r.group.PUT(path, bunrouter.HTTPHandlerFunc(handler))
			case DELETE:
				r.group.DELETE(path, bunrouter.HTTPHandlerFunc(handler))
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		}
		return r
	}
}

func WithRouterProducerHandler(method string, path string, bodyMap stateless.OneToOneFunction[structure.Bytes, structure.Bytes, structure.Bytes, structure.Bytes]) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		handlerFunction := router.NewRouteProducer(
			router.WithRouteProducerRuntime(r.producer),
			router.WithRouteBodyMap(bodyMap),
		)
		if r.group == nil {
			switch strings.ToUpper(method) {
			case GET:
				r.router.GET(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case POST:
				r.router.POST(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case PUT:
				r.router.PUT(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case DELETE:
				r.router.DELETE(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		} else {
			switch strings.ToUpper(method) {
			case GET:
				r.group.GET(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case POST:
				r.group.POST(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case PUT:
				r.group.PUT(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			case DELETE:
				r.group.DELETE(path, bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
			default:
				logger.Warn("unknown method", zap.String("method", method))
			}
		}
		return r
	}
}

func WithRouterFlow(configurations ...runtime.Configuration[*router.RouteFlow]) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.flowConfiguration = append(r.flowConfiguration, configurations...)
		return r
	}
}

func WithRouterPrometheus() runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.router.GET("/prometheus", bunrouter.HTTPHandlerFunc(promhttp.Handler().ServeHTTP))
		return r
	}
}

func WithRouterProducer(producer flow.Producer) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.producer = producer
		return r
	}
}

// implementation
type Router struct {
	runtimeConfiguration []runtime.Configuration[*runtime.HttpRunnable]
	flowConfiguration    []runtime.Configuration[*router.RouteFlow]
	router               *bunrouter.Router
	group                *bunrouter.Group
	producer             flow.Producer
}
