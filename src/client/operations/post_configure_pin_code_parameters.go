// Code generated by go-swagger; DO NOT EDIT.

package operations

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

	models "github.com/skycoin/hardware-wallet-daemon/src/models"
)

// NewPostConfigurePinCodeParams creates a new PostConfigurePinCodeParams object
// with the default values initialized.
func NewPostConfigurePinCodeParams() *PostConfigurePinCodeParams {
	var ()
	return &PostConfigurePinCodeParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostConfigurePinCodeParamsWithTimeout creates a new PostConfigurePinCodeParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostConfigurePinCodeParamsWithTimeout(timeout time.Duration) *PostConfigurePinCodeParams {
	var ()
	return &PostConfigurePinCodeParams{

		timeout: timeout,
	}
}

// NewPostConfigurePinCodeParamsWithContext creates a new PostConfigurePinCodeParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostConfigurePinCodeParamsWithContext(ctx context.Context) *PostConfigurePinCodeParams {
	var ()
	return &PostConfigurePinCodeParams{

		Context: ctx,
	}
}

// NewPostConfigurePinCodeParamsWithHTTPClient creates a new PostConfigurePinCodeParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostConfigurePinCodeParamsWithHTTPClient(client *http.Client) *PostConfigurePinCodeParams {
	var ()
	return &PostConfigurePinCodeParams{
		HTTPClient: client,
	}
}

/*PostConfigurePinCodeParams contains all the parameters to send to the API endpoint
for the post configure pin code operation typically these are written to a http.Request
*/
type PostConfigurePinCodeParams struct {

	/*ConfigurePinCodeRequest
	  ConfigurePinCodeRequest is request data for /api/v1/configure_pin_code

	*/
	ConfigurePinCodeRequest *models.ConfigurePinCodeRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post configure pin code params
func (o *PostConfigurePinCodeParams) WithTimeout(timeout time.Duration) *PostConfigurePinCodeParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post configure pin code params
func (o *PostConfigurePinCodeParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post configure pin code params
func (o *PostConfigurePinCodeParams) WithContext(ctx context.Context) *PostConfigurePinCodeParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post configure pin code params
func (o *PostConfigurePinCodeParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post configure pin code params
func (o *PostConfigurePinCodeParams) WithHTTPClient(client *http.Client) *PostConfigurePinCodeParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post configure pin code params
func (o *PostConfigurePinCodeParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithConfigurePinCodeRequest adds the configurePinCodeRequest to the post configure pin code params
func (o *PostConfigurePinCodeParams) WithConfigurePinCodeRequest(configurePinCodeRequest *models.ConfigurePinCodeRequest) *PostConfigurePinCodeParams {
	o.SetConfigurePinCodeRequest(configurePinCodeRequest)
	return o
}

// SetConfigurePinCodeRequest adds the configurePinCodeRequest to the post configure pin code params
func (o *PostConfigurePinCodeParams) SetConfigurePinCodeRequest(configurePinCodeRequest *models.ConfigurePinCodeRequest) {
	o.ConfigurePinCodeRequest = configurePinCodeRequest
}

// WriteToRequest writes these params to a swagger request
func (o *PostConfigurePinCodeParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ConfigurePinCodeRequest != nil {
		if err := r.SetBodyParam(o.ConfigurePinCodeRequest); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
