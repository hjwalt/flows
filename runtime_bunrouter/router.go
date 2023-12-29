package runtime_bunrouter

import (
	"context"
	"net/http"

	"github.com/hjwalt/routes/runtime_chi"
	"github.com/hjwalt/runway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// THIS HAS BEEN RECODED TO BE A FACADE TO routes REPOSITORY
// TODO: complete removal

// constructor
func NewRouter(configurations ...runtime.Configuration[*Router]) runtime.Runtime {
	return runtime_chi.New[context.Context](configurations...)
}

// configuration
func WithRouterPort(port int) runtime.Configuration[*Router] {
	return runtime_chi.WithPort[context.Context](port)
}

func WithRouterHttpConfiguration(config runtime.Configuration[*runtime.HttpRunnable]) runtime.Configuration[*Router] {
	return runtime_chi.WithHttpConfiguration[context.Context](config)
}

func WithRouterHttpHandler(method string, path string, handler http.HandlerFunc) runtime.Configuration[*Router] {
	return runtime_chi.WithHttpHandler[context.Context](path, method, handler)
}

func WithRouterPrometheus() runtime.Configuration[*Router] {
	return runtime_chi.WithHttpHandler[context.Context]("/prometheus", "GET", promhttp.Handler())
}

// implementation
type Router = runtime_chi.Runtime[context.Context]
