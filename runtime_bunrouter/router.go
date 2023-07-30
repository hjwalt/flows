package runtime_bunrouter

import (
	"net/http"
	"strings"

	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/stateless"
	"github.com/hjwalt/runway/logger"
	"github.com/hjwalt/runway/runtime"
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

// constructor
func NewRouter(configurations ...runtime.Configuration[*Router]) runtime.Runtime {
	router := &Router{
		runtimeConfiguration: []runtime.Configuration[*runtime.HttpRunnable]{},
		router:               bunrouter.New(),
	}

	for _, configuration := range configurations {
		router = configuration(router)
	}

	runtimeConfiguration := append(router.runtimeConfiguration, runtime.HttpWithHandler(router.router))

	return runtime.NewHttp(runtimeConfiguration...)
}

// configuration
func WithRouterPort(port int) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
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

func WithRouterProducerHandler(method string, path string, bodyMap stateless.OneToOneFunction[message.Bytes, message.Bytes, message.Bytes, message.Bytes]) runtime.Configuration[*Router] {
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
		handlerFunction := router.NewRouteFlow(
			configurations...,
		)
		if r.group == nil {
			r.router.GET("/flow", bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
		} else {
			r.group.GET("/flow", bunrouter.HTTPHandlerFunc(handlerFunction.Handle))
		}
		return r
	}
}

func WithRouterPrometheus() runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		if r.group == nil {
			r.router.GET("/prometheus", bunrouter.HTTPHandlerFunc(promhttp.Handler().ServeHTTP))
		} else {
			r.group.GET("/prometheus", bunrouter.HTTPHandlerFunc(promhttp.Handler().ServeHTTP))
		}
		return r
	}
}

func WithRouterProducer(producer message.Producer) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.producer = producer
		return r
	}
}

// implementation
type Router struct {
	runtimeConfiguration []runtime.Configuration[*runtime.HttpRunnable]
	router               *bunrouter.Router
	group                *bunrouter.Group
	producer             message.Producer
}
