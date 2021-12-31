package pipekit

import (
	"net/http"
	"net/url"
)

// Params holds query parameters
type Params struct {
	Values url.Values
}

// ParamsContainer allows backend invokers to configure parameters
type ParamsContainer interface {
	GetParams() *Params
}

// ResponseSetter allows backend invokers to configure the api response
type ResponseSetter interface {
	SetResponse(response *APIResponse)
}

// APIResponse houses the standard golang http response plus some added info
type APIResponse struct {
	http.Response
}
