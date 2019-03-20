// Code generated by go-swagger; DO NOT EDIT.

package nodes

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// CreateNodeCreatedCode is the HTTP code returned for type CreateNodeCreated
const CreateNodeCreatedCode int = 201

/*CreateNodeCreated Created

swagger:response createNodeCreated
*/
type CreateNodeCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Node `json:"body,omitempty"`
}

// NewCreateNodeCreated creates CreateNodeCreated with default headers values
func NewCreateNodeCreated() *CreateNodeCreated {

	return &CreateNodeCreated{}
}

// WithPayload adds the payload to the create node created response
func (o *CreateNodeCreated) WithPayload(payload *models.Node) *CreateNodeCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create node created response
func (o *CreateNodeCreated) SetPayload(payload *models.Node) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNodeCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateNodeDefault error

swagger:response createNodeDefault
*/
type CreateNodeDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateNodeDefault creates CreateNodeDefault with default headers values
func NewCreateNodeDefault(code int) *CreateNodeDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateNodeDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create node default response
func (o *CreateNodeDefault) WithStatusCode(code int) *CreateNodeDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create node default response
func (o *CreateNodeDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create node default response
func (o *CreateNodeDefault) WithPayload(payload *models.Error) *CreateNodeDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create node default response
func (o *CreateNodeDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNodeDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}