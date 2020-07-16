// Code generated by go-swagger; DO NOT EDIT.

package configuration

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
)

// NewGetConfigBridgeParams creates a new GetConfigBridgeParams object
// no default values defined in spec.
func NewGetConfigBridgeParams() GetConfigBridgeParams {

	return GetConfigBridgeParams{}
}

// GetConfigBridgeParams contains all the bound params for the get config bridge operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetConfigBridge
type GetConfigBridgeParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetConfigBridgeParams() beforehand.
func (o *GetConfigBridgeParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}