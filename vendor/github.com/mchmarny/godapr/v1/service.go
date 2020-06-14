package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// InvokeServiceWithData invokes the remote service method
func (c *Client) InvokeServiceWithData(ctx trace.SpanContext, service, method string, in []byte) (out []byte, err error) {
	if service == "" {
		return nil, errors.New("nil service")
	}
	if method == "" {
		return nil, errors.New("nil method")
	}
	url := fmt.Sprintf("%s/v1.0/invoke/%s/method/%s", c.url, service, method)

	var content *bytes.Buffer
	if in != nil {
		content = bytes.NewBuffer(in)
	}

	req, err := http.NewRequest(http.MethodPost, url, content)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating invoking request: %s", url)
	}

	out, status, err := c.exec(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "error executing: %+v", req)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("invalid response code to %s: %d", url, status)
	}

	return out, nil
}

// InvokeServiceWithIdentity serializes input data to JSON and invokes InvokeServiceWithData
func (c *Client) InvokeServiceWithIdentity(ctx trace.SpanContext, service, method string, in interface{}) (out []byte, err error) {
	if in == nil {
		return nil, errors.New("nil input identity")
	}
	b, err := json.Marshal(in)
	if err != nil {
		return nil, errors.Wrapf(err, "error serializing identity: %v", in)
	}
	return c.InvokeServiceWithData(ctx, service, method, b)
}

// InvokeService and invokes InvokeServiceWithData without any payload
func (c *Client) InvokeService(ctx trace.SpanContext, service, method string) (out []byte, err error) {
	return c.InvokeServiceWithData(ctx, service, method, nil)
}
