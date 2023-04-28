package router

import "errors"

// route produce

var ErrorRouteGettingRequestBody = errors.New("route getting request body")
var ErrorRouteReadingRequestBody = errors.New("route reading request body")
var ErrorRouteParsingRequestBody = errors.New("route parsing request body")
var ErrorRouteMappingRequestBody = errors.New("route mapping request body")
var ErrorRouteProducingMessage = errors.New("route producing message")

// write json

var ErrorWriteJsonFormat = errors.New("write json formatting")
