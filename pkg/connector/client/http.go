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
)

type customTransport struct {
	userAgent string
	base      http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clonedReq := req.Clone(req.Context())
	clonedReq.Header.Set("User-Agent", t.userAgent)

	return t.base.RoundTrip(clonedReq)
}

func NewSGNLHttpClient(timeout time.Duration, userAgent string) *http.Client {
	// This is a default value if a SGNL 1P adapter does not define a custom user agent.
	if userAgent == "" {
		userAgent = "sgnl-adapter"
	}

	t := &customTransport{
		base:      http.DefaultTransport,
		userAgent: userAgent,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: t,
	}
}

// NewSGNLHTTPClientWithProxy for proxying http requests based on the Connector context
// present in the request.
// TODO: Replace `NewSGNLHttpClient()` with `NewSGNLHTTPClientWithProxy` for adapters with protocols suppurted by
// the `ProxyServiceClient“.
func NewSGNLHTTPClientWithProxy(timeout time.Duration, userAgent string, client v1proxy.ProxyServiceClient) *http.Client {
	// This is a default value if a SGNL 1P adapter does not define a custom user agent.
	if userAgent == "" {
		userAgent = "sgnl-adapter"
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
	_, ok := connector.FromContext(ctx)

	if !ok {
		return t.rt.RoundTrip(req)
	}

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

	// Create proxied HTTP request message.
	// TODO: to include connector metadata in the request.
	grpcReq := &v1proxy.Request{
		RequestType: &v1proxy.Request_HttpRequest{
			HttpRequest: &v1proxy.HTTPRequest{
				Method: req.Method,
				Url:    req.URL.String(),
				Headers: func() map[string]*v1proxy.StringValues {
					headers := make(map[string]*v1proxy.StringValues)
					for k, v := range req.Header {
						if len(v) > 0 {
							headers[k] = &v1proxy.StringValues{
								Values: v,
							}
						}
					}

					return headers
				}(),
				Body: body,
			},
		},
	}

	// Send request.
	resp, err := t.proxyClient.ProxyRequest(req.Context(), grpcReq)
	if err != nil {
		return nil, err
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
