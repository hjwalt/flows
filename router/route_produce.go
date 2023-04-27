package router

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/hjwalt/flows/format"
	"github.com/hjwalt/flows/message"
	"github.com/hjwalt/flows/runtime"
	"github.com/uptrace/bunrouter"
)

const (
	RouteHeaderPrefix     = "HTTP_HEADER_"
	RouteParamPrefix      = "HTTP_PARAM_"
	RouteRequestKeyHeader = "HTTP_HEADER_X_REQUEST_KEY"
)

// route

type RouteProduceMapFunction func(ctx context.Context, req message.Message[message.Bytes, message.Bytes]) (message.Message[message.Bytes, message.Bytes], error)

type RouteProducerResponse struct {
	Message string
}

// constructor
func NewRouteProducer(configurations ...runtime.Configuration[*RouteProducer]) *RouteProducer {
	rp := &RouteProducer{}
	for _, configuration := range configurations {
		rp = configuration(rp)
	}
	return rp
}

// configurations
func WithRouteProducerRuntime(producer runtime.Producer) runtime.Configuration[*RouteProducer] {
	return func(psf *RouteProducer) *RouteProducer {
		psf.producer = producer
		return psf
	}
}

func WithRouteBodyMap(bodyMap RouteProduceMapFunction) runtime.Configuration[*RouteProducer] {
	return func(psf *RouteProducer) *RouteProducer {
		psf.bodyMap = bodyMap
		return psf
	}
}

// implementation
type RouteProducer struct {
	producer runtime.Producer
	bodyMap  RouteProduceMapFunction
}

func (rp *RouteProducer) Route(w http.ResponseWriter, req bunrouter.Request) error {
	requestHeaders := map[string][]message.Bytes{}

	// request headers
	for hKey, hValues := range req.Header {
		msgHeaderValues := make([]message.Bytes, len(hValues))
		for hIndex, hValue := range hValues {
			msgHeaderValues[hIndex] = []byte(hValue)
		}

		msgHeaderKey := RouteHeaderPrefix + processHeaderKey(hKey)
		requestHeaders[msgHeaderKey] = msgHeaderValues
	}

	// request params
	for hKey, hValues := range req.URL.Query() {
		msgHeaderValues := make([]message.Bytes, len(hValues))
		for hIndex, hValue := range hValues {
			msgHeaderValues[hIndex] = []byte(hValue)
		}

		msgHeaderKey := RouteParamPrefix + processHeaderKey(hKey)
		requestHeaders[msgHeaderKey] = msgHeaderValues
	}

	// request key
	requestKey := ""
	if requestKeyHeader, requestKeyHeaderExist := requestHeaders[RouteRequestKeyHeader]; requestKeyHeaderExist {
		requestKey = string(requestKeyHeader[0])
	}

	// request body
	defer req.Body.Close()
	body, readErr := io.ReadAll(req.Body)
	if readErr != nil {
		return errors.Join(ErrorRouteReadingRequestBody, readErr)
	}

	requestMessage := message.Message[message.Bytes, message.Bytes]{
		Key:     []byte(requestKey),
		Value:   body,
		Headers: requestHeaders,
	}

	// mapping body
	messageMessage, requestMapErr := rp.bodyMap(req.Context(), requestMessage)
	if requestMapErr != nil {
		return errors.Join(ErrorRouteMappingRequestBody, requestMapErr)
	}

	produceErr := rp.producer.Produce(req.Context(), []message.Message[message.Bytes, message.Bytes]{messageMessage})
	if produceErr != nil {
		return errors.Join(ErrorRouteProducingMessage, produceErr)
	}

	return WriteJson(w, 200, RouteProducerResponse{Message: "ok"}, format.Json[RouteProducerResponse]())
}

func processHeaderKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, "-", "_")
	return key
}
