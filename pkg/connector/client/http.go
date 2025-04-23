// Copyright 2025 SGNL.ai, Inc.

package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	v1proxy "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
	"google.golang.org/grpc/metadata"
)

// NewSGNLHTTPClientWithProxy for proxying http requests based on the Connector context
// present in the request.
func NewSGNLHTTPClientWithProxy(timeout time.Duration, userAgent string, client v1proxy.ProxyServiceClient) *http.Client {
	// This is a default value if custom user agent is not specified.
	if userAgent == "" {
		userAgent = "sgnl-"
	}

	return &http.Client{
		Timeout: timeout,
		Transport: &Transport{
			rt:              http.DefaultTransport,
			userAgentHeader: userAgent,
			proxyClient:     client,
		},
	}
}

// Transport with proxy client.
type Transport struct {
	rt              http.RoundTripper
	userAgentHeader string
	proxyClient     v1proxy.ProxyServiceClient
}

// RoundTrip implementation to either call the underlying default HTTP transport or
// proxy the request to the connector service using the proxy client.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// According to the RoundTripper spec -
	// RoundTrip should not modify the request, except for
	// consuming and closing the Request's Body.
	req = req.Clone(req.Context()) // According to RoundTripper spec, we shouldn't modify the original request.
	req.Header.Set("User-Agent", t.userAgentHeader)

	// Lookup Connector context for proxying the request.
	// In case, context is not present the request is forward using default
	// HTTP round-tripper.
	ctx := req.Context()
	ci, ok := connector.FromContext(ctx)

	if !ok {
		return t.rt.RoundTrip(req)
	}

	// Update the metadata context with connector information
	ctx = metadata.AppendToOutgoingContext(req.Context(),
		connector.METADATA_CLIENT_ID, ci.ClientID,
		connector.METADATA_CONNECTOR_ID, ci.ID,
		connector.METADATA_TENANT_ID, ci.TenantID)

	// Prepare the request payload for the proxied request.
	var body []byte

	if req.Body != nil {
		var err error

		if body, err = io.ReadAll(req.Body); err != nil {
			return nil, err
		}

		req.Body.Close()
		// Reset the request body for potential retries.
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	// Copy all the header key and values from the request header.
	headers := make(map[string]*v1proxy.StringValues)
	for k, v := range req.Header {
		if len(v) > 0 {
			headers[k] = &v1proxy.StringValues{
				Values: v,
			}
		}
	}

	// Create proxied HTTP request message.
	// TODO: to include connector metadata in the request.
	grpcReq := &v1proxy.Request{
		RequestType: &v1proxy.Request_HttpRequest{
			HttpRequest: &v1proxy.HTTPRequest{
				Method:  req.Method,
				Url:     req.URL.String(),
				Headers: headers,
				Body:    body,
			},
		},
	}

	// Send request.
	resp, err := t.proxyClient.ProxyRequest(ctx, grpcReq)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("%s", resp.Error)
	}

	// Convert gRPC response message to the http.Response
	httpResp := &http.Response{
		StatusCode: int(resp.GetHttpResponse().StatusCode),
		Status:     fmt.Sprintf("%d %s", resp.GetHttpResponse().StatusCode, http.StatusText(int(resp.GetHttpResponse().StatusCode))),
		Body:       io.NopCloser(bytes.NewReader(resp.GetHttpResponse().Body)),
		Header:     make(http.Header),
		Request:    req,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	// Add response headers
	for k, v := range resp.GetHttpResponse().Headers {
		if v != nil {
			for _, val := range v.GetValues() {
				httpResp.Header.Add(k, val)
			}
		}
	}

	return httpResp, nil
}
