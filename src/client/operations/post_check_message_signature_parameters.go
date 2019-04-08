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

// NewPostCheckMessageSignatureParams creates a new PostCheckMessageSignatureParams object
// with the default values initialized.
func NewPostCheckMessageSignatureParams() *PostCheckMessageSignatureParams {
	var ()
	return &PostCheckMessageSignatureParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostCheckMessageSignatureParamsWithTimeout creates a new PostCheckMessageSignatureParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostCheckMessageSignatureParamsWithTimeout(timeout time.Duration) *PostCheckMessageSignatureParams {
	var ()
	return &PostCheckMessageSignatureParams{

		timeout: timeout,
	}
}

// NewPostCheckMessageSignatureParamsWithContext creates a new PostCheckMessageSignatureParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostCheckMessageSignatureParamsWithContext(ctx context.Context) *PostCheckMessageSignatureParams {
	var ()
	return &PostCheckMessageSignatureParams{

		Context: ctx,
	}
}

// NewPostCheckMessageSignatureParamsWithHTTPClient creates a new PostCheckMessageSignatureParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostCheckMessageSignatureParamsWithHTTPClient(client *http.Client) *PostCheckMessageSignatureParams {
	var ()
	return &PostCheckMessageSignatureParams{
		HTTPClient: client,
	}
}

/*PostCheckMessageSignatureParams contains all the parameters to send to the API endpoint
for the post check message signature operation typically these are written to a http.Request
*/
type PostCheckMessageSignatureParams struct {

	/*CheckMessageSignatureRequest
	  CheckMessageSignatureRequest is request data for /api/checkMessageSignature

	*/
	CheckMessageSignatureRequest *models.CheckMessageSignatureRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the post check message signature params
func (o *PostCheckMessageSignatureParams) WithTimeout(timeout time.Duration) *PostCheckMessageSignatureParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post check message signature params
func (o *PostCheckMessageSignatureParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post check message signature params
func (o *PostCheckMessageSignatureParams) WithContext(ctx context.Context) *PostCheckMessageSignatureParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post check message signature params
func (o *PostCheckMessageSignatureParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post check message signature params
func (o *PostCheckMessageSignatureParams) WithHTTPClient(client *http.Client) *PostCheckMessageSignatureParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post check message signature params
func (o *PostCheckMessageSignatureParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCheckMessageSignatureRequest adds the checkMessageSignatureRequest to the post check message signature params
func (o *PostCheckMessageSignatureParams) WithCheckMessageSignatureRequest(checkMessageSignatureRequest *models.CheckMessageSignatureRequest) *PostCheckMessageSignatureParams {
	o.SetCheckMessageSignatureRequest(checkMessageSignatureRequest)
	return o
}

// SetCheckMessageSignatureRequest adds the checkMessageSignatureRequest to the post check message signature params
func (o *PostCheckMessageSignatureParams) SetCheckMessageSignatureRequest(checkMessageSignatureRequest *models.CheckMessageSignatureRequest) {
	o.CheckMessageSignatureRequest = checkMessageSignatureRequest
}

// WriteToRequest writes these params to a swagger request
func (o *PostCheckMessageSignatureParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.CheckMessageSignatureRequest != nil {
		if err := r.SetBodyParam(o.CheckMessageSignatureRequest); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}