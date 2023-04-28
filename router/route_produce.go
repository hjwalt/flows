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
	"github.com/hjwalt/runway/logger"
)

const (
	RouteHeaderPrefix     = "HTTP_HEADER_"
	RouteParamPrefix      = "HTTP_PARAM_"
	RouteRequestKeyHeader = "HTTP_HEADER_X_REQUEST_KEY"
)

// route

type BasicResponse struct {
	Message string `json:"message"`
}

type RouteProduceMapFunction func(ctx context.Context, req message.Message[message.Bytes, message.Bytes]) (message.Message[message.Bytes, message.Bytes], error)

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

func (rp *RouteProducer) Handle(w http.ResponseWriter, req *http.Request) {
	err := rp.Produce(w, req)
	if err != nil {
		logger.ErrorErr("handle error", err)
		WriteError(w, 500, errors.New("server error"))
	}
}

func (rp *RouteProducer) Produce(w http.ResponseWriter, req *http.Request) error {
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
	var requestBody []byte
	if req.Body != nil {
		defer req.Body.Close()
		body, readErr := io.ReadAll(req.Body)
		if readErr != nil {
			return errors.Join(ErrorRouteReadingRequestBody, readErr)
		}

		requestBody = body
	} else {
		requestBody = make([]byte, 0)
	}

	requestMessage := message.Message[message.Bytes, message.Bytes]{
		Key:     []byte(requestKey),
		Value:   requestBody,
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

	return WriteJson(w, 200, BasicResponse{Message: "ok"}, format.Json[BasicResponse]())
}

func processHeaderKey(key string) string {
	key = strings.ToUpper(key)
	key = strings.ReplaceAll(key, "-", "_")
	return key
}
