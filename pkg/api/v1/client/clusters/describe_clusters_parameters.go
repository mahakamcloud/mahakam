// Code generated by go-swagger; DO NOT EDIT.

package clusters

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewDescribeClustersParams creates a new DescribeClustersParams object
// with the default values initialized.
func NewDescribeClustersParams() *DescribeClustersParams {
	var ()
	return &DescribeClustersParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewDescribeClustersParamsWithTimeout creates a new DescribeClustersParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewDescribeClustersParamsWithTimeout(timeout time.Duration) *DescribeClustersParams {
	var ()
	return &DescribeClustersParams{

		timeout: timeout,
	}
}

// NewDescribeClustersParamsWithContext creates a new DescribeClustersParams object
// with the default values initialized, and the ability to set a context for a request
func NewDescribeClustersParamsWithContext(ctx context.Context) *DescribeClustersParams {
	var ()
	return &DescribeClustersParams{

		Context: ctx,
	}
}

// NewDescribeClustersParamsWithHTTPClient creates a new DescribeClustersParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewDescribeClustersParamsWithHTTPClient(client *http.Client) *DescribeClustersParams {
	var ()
	return &DescribeClustersParams{
		HTTPClient: client,
	}
}

/*DescribeClustersParams contains all the parameters to send to the API endpoint
for the describe clusters operation typically these are written to a http.Request
*/
type DescribeClustersParams struct {

	/*Name
	  Cluster name

	*/
	Name *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the describe clusters params
func (o *DescribeClustersParams) WithTimeout(timeout time.Duration) *DescribeClustersParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the describe clusters params
func (o *DescribeClustersParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the describe clusters params
func (o *DescribeClustersParams) WithContext(ctx context.Context) *DescribeClustersParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the describe clusters params
func (o *DescribeClustersParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the describe clusters params
func (o *DescribeClustersParams) WithHTTPClient(client *http.Client) *DescribeClustersParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the describe clusters params
func (o *DescribeClustersParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithName adds the name to the describe clusters params
func (o *DescribeClustersParams) WithName(name *string) *DescribeClustersParams {
	o.SetName(name)
	return o
}

// SetName adds the name to the describe clusters params
func (o *DescribeClustersParams) SetName(name *string) {
	o.Name = name
}

// WriteToRequest writes these params to a swagger request
func (o *DescribeClustersParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Name != nil {

		// query param name
		var qrName string
		if o.Name != nil {
			qrName = *o.Name
		}
		qName := qrName
		if qName != "" {
			if err := r.SetQueryParam("name", qName); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
