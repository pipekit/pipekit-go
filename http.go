package pipekit

import (
	"net/url"
)

// Params holds query parameters
type Params struct {
	Values url.Values
}

// GetParams returns a params struct so that any structs that inherit params
// will conform to the ParamsContainer interface
func (p *Params) GetParams() *Params {
	return p
}

// ParamsContainer allows backend invokers to configure parameters
type ParamsContainer interface {
	GetParams() *Params
}
