// Code generated by go-swagger; DO NOT EDIT.

package gre_networks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// GetGreNetworksOKCode is the HTTP code returned for type GetGreNetworksOK
const GetGreNetworksOKCode int = 200

/*GetGreNetworksOK list created GRE networks

swagger:response getGreNetworksOK
*/
type GetGreNetworksOK struct {

	/*
	  In: Body
	*/
	Payload []*models.GreNetwork `json:"body,omitempty"`
}

// NewGetGreNetworksOK creates GetGreNetworksOK with default headers values
func NewGetGreNetworksOK() *GetGreNetworksOK {

	return &GetGreNetworksOK{}
}

// WithPayload adds the payload to the get gre networks o k response
func (o *GetGreNetworksOK) WithPayload(payload []*models.GreNetwork) *GetGreNetworksOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get gre networks o k response
func (o *GetGreNetworksOK) SetPayload(payload []*models.GreNetwork) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetGreNetworksOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.GreNetwork, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GetGreNetworksDefault generic error response

swagger:response getGreNetworksDefault
*/
type GetGreNetworksDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetGreNetworksDefault creates GetGreNetworksDefault with default headers values
func NewGetGreNetworksDefault(code int) *GetGreNetworksDefault {
	if code <= 0 {
		code = 500
	}

	return &GetGreNetworksDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get gre networks default response
func (o *GetGreNetworksDefault) WithStatusCode(code int) *GetGreNetworksDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get gre networks default response
func (o *GetGreNetworksDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get gre networks default response
func (o *GetGreNetworksDefault) WithPayload(payload *models.Error) *GetGreNetworksDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get gre networks default response
func (o *GetGreNetworksDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetGreNetworksDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
