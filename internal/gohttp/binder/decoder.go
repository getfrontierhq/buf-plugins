package binder

import "net/http"

type RequestDecoder struct {
	Request *http.Request
}

type ResponseDecoder struct {
	Response *http.Response
}
