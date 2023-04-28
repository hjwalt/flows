package runtime_bunrouter

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hjwalt/flows/router"
	"github.com/hjwalt/flows/runtime"
	"github.com/hjwalt/runway/logger"
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
func NewRouter(configurations ...runtime.Configuration[*Router]) *Router {
	router := &Router{
		router: bunrouter.New(),
	}
	for _, configuration := range configurations {
		router = configuration(router)
	}
	return router
}

// configuration
func WithRouterPort(port int) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.port = port
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

func WithRouterProducerHandler(method string, path string, bodyMap router.RouteProduceMapFunction) runtime.Configuration[*Router] {
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

func WithRouterProducer(producer runtime.Producer) runtime.Configuration[*Router] {
	return func(r *Router) *Router {
		r.producer = producer
		return r
	}
}

// implementation
type Router struct {
	server   *http.Server
	router   *bunrouter.Router
	group    *bunrouter.Group
	port     int
	producer runtime.Producer
}

func (r *Router) Start() error {
	if r.port == 0 {
		r.port = 8080
	}

	// run server
	r.server = &http.Server{
		Addr:         ":" + fmt.Sprintf("%d", r.port),
		Handler:      r.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go r.Run()

	logger.Infof("router started")

	return nil
}

func (runtime *Router) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := runtime.server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: ", zap.Error(err))
	}
	logger.Infof("router stopped")
}

func (r *Router) Run() {
	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	logger.Infof("router run ended")
}
