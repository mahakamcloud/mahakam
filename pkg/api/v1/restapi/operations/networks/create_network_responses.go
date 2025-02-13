// Code generated by go-swagger; DO NOT EDIT.

package networks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/mahakamcloud/mahakam/pkg/api/v1/models"
)

// CreateNetworkCreatedCode is the HTTP code returned for type CreateNetworkCreated
const CreateNetworkCreatedCode int = 201

/*CreateNetworkCreated Created

swagger:response createNetworkCreated
*/
type CreateNetworkCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Network `json:"body,omitempty"`
}

// NewCreateNetworkCreated creates CreateNetworkCreated with default headers values
func NewCreateNetworkCreated() *CreateNetworkCreated {

	return &CreateNetworkCreated{}
}

// WithPayload adds the payload to the create network created response
func (o *CreateNetworkCreated) WithPayload(payload *models.Network) *CreateNetworkCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create network created response
func (o *CreateNetworkCreated) SetPayload(payload *models.Network) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNetworkCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateNetworkDefault error

swagger:response createNetworkDefault
*/
type CreateNetworkDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateNetworkDefault creates CreateNetworkDefault with default headers values
func NewCreateNetworkDefault(code int) *CreateNetworkDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateNetworkDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create network default response
func (o *CreateNetworkDefault) WithStatusCode(code int) *CreateNetworkDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create network default response
func (o *CreateNetworkDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create network default response
func (o *CreateNetworkDefault) WithPayload(payload *models.Error) *CreateNetworkDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create network default response
func (o *CreateNetworkDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNetworkDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
