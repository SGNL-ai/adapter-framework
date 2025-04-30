// Copyright 2025 SGNL.ai, Inc.

package client

import (
	context "context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	v1proxy "github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1"
)

func TestCustomUserAgent(t *testing.T) {
	// Expected User-Agent value
	tests := map[string]struct {
		inputUserAgent string
		wantUserAgent  string
	}{
		"empty input": {
			inputUserAgent: "",
			wantUserAgent:  "sgnl-",
		},
		"non-empty input": {
			inputUserAgent: "sgnl-myAdapter/1.0.0",
			wantUserAgent:  "sgnl-myAdapter/1.0.0",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a test server to capture the request
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check the User-Agent header
				if gotUserAgent := r.Header.Get("User-Agent"); gotUserAgent != tc.wantUserAgent {
					w.WriteHeader(http.StatusBadRequest)
					t.Fatalf("unexpected User-Agent header: got %q, want %q", gotUserAgent, tc.wantUserAgent)

					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer testServer.Close()

			// Create an HTTP client with the custom User-Agent
			client := NewSGNLHTTPClientWithProxy(time.Second, tc.inputUserAgent, nil)

			// Fire a request
			resp, err := client.Get(testServer.URL)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Ensure the response status is OK
			if resp.StatusCode != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", resp.StatusCode, http.StatusOK)
			}
		})
	}
}

func TestGivenSGNLHTTPClientWithTimeoutThenRequestTimesout(t *testing.T) {
	// Build
	block := make(chan struct{})

	// Start test server with infinite delay
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block // Block indefinitely until test completes
		w.Write([]byte("OK"))
	}))
	defer ts.Close()
	defer close(block) // Ensure channel cleanup

	client := NewSGNLHTTPClientWithProxy(time.Second, "", nil)

	// Test
	_, err := client.Get(ts.URL)
	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}

	// Verify
	if !os.IsTimeout(err) {
		t.Errorf("Expected timeout error, got: %v", err)
	}
}

func TestGivenSGNLHTTPClientWithoutGRPCProxyThenSendRequestSentWithoutProxyAndReturnsStatusOK(t *testing.T) {
	// Build
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewSGNLHTTPClientWithProxy(time.Second, "", nil)

	// Test
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal("Failed to send a request to the Test server using SGNL client")
	}

	// Verify
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test server failed to return StatusOK using SGNL client")
	}
}

func TestGivenSGNLHTTPClientWithConnectorContextAndWithoutGRPCProxyThenSendRequestSentWithoutProxyAndReturnsStatusOK(t *testing.T) {
	// Build
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewSGNLHTTPClientWithProxy(time.Second, "", nil)

	ctx, err := connector.WithContext(context.Background(), connector.ConnectorInfo{})
	if err != nil {
		t.Fatalf("failed to create connector info context %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create a http request, %v", err)
	}

	// Test
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send a request to the Test server using SGNL client, %v", err)
	}
	defer resp.Body.Close()

	// Verify
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test server failed to return StatusOK using SGNL client")
	}
}

var GRPCTestServerResponse = "response from the grpc test server"

type testServer struct {
	ci *connector.ConnectorInfo
	v1proxy.UnimplementedProxyServiceServer
}

func (s *testServer) ProxyRequest(ctx context.Context, req *v1proxy.ProxyRequestMessage) (*v1proxy.Response, error) {
	if s.ci != nil {
		if req.ClientId != s.ci.ClientID {
			return nil, fmt.Errorf("Expected %v, got %v client id", req.ClientId, s.ci.ClientID)
		} else if req.ConnectorId != s.ci.ID {
			return nil, fmt.Errorf("Expected %v, got %v connector id", req.ConnectorId, s.ci.ID)
		} else if req.TenantId != s.ci.TenantID {
			return nil, fmt.Errorf("Expected %v, got %v tenant id", req.TenantId, s.ci.TenantID)
		}
	}

	return &v1proxy.Response{
		ResponseType: &v1proxy.Response_HttpResponse{
			HttpResponse: &v1proxy.HTTPResponse{
				StatusCode: http.StatusOK,
				Body:       []byte(GRPCTestServerResponse),
			}}}, nil
}

func TestGivenSGNLHTTPClientWithGRPCProxyAndRequestWithoutConnectorContextThenRequestSentWithoutProxyAndReturnsStatusOK(t *testing.T) {
	// Build
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create in-memory listener
	lis := bufconn.Listen(1024 * 1024)
	defer lis.Close()

	// Create gRPC test server and serve from the in-memory listener
	srv := grpc.NewServer()
	v1proxy.RegisterProxyServiceServer(srv, &testServer{})
	go srv.Serve(lis)
	defer srv.Stop()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := NewSGNLHTTPClientWithProxy(time.Second, "", v1proxy.NewProxyServiceClient(conn))

	// Test
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("Failed to send a request to the Test server using SGNL client, %v", err)
	}

	// Verify
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test server failed to return StatusOK using SGNL client")
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if len(body) > 0 {
		t.Errorf("Expected proxy response %s got %s", GRPCTestServerResponse, string(body))
	}
}

func TestGivenSGNLHTTPClientWithGRPCProxyAndRequestWithConnectorContextThenRequestSentToProxyAndReturnsStatusOK(t *testing.T) {
	// Build
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Create in-memory listener
	lis := bufconn.Listen(1024 * 1024)
	defer lis.Close()

	connectorInfo := connector.ConnectorInfo{
		ID:         "123-456-789",
		ClientID:   "client-id",
		SourceType: connector.Datasource,
		SourceID:   "datasource-id",
	}

	// Create gRPC test server and serve from the in-memory listener
	srv := grpc.NewServer()
	v1proxy.RegisterProxyServiceServer(srv, &testServer{
		ci: &connectorInfo,
	})
	go srv.Serve(lis)
	defer srv.Stop()

	conn, err := grpc.DialContext(context.Background(), v1proxy.ProxyService_ServiceDesc.ServiceName,
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := NewSGNLHTTPClientWithProxy(time.Second, "", v1proxy.NewProxyServiceClient(conn))

	// Create and pass the connector context to the http request.
	ctx, err := connector.WithContext(context.Background(), connectorInfo)
	if err != nil {
		t.Fatalf("failed to create connector info context %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf("failed to create a http request, %v", err)
	}

	// Test
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send a request to the Test server using SGNL client, %v", err)
	}

	// Verify
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Test proxy server failed to return StatusOK using SGNL client")
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if GRPCTestServerResponse != string(body) {
		t.Errorf("Expected proxy response %s got %s", GRPCTestServerResponse, string(body))
	}
}
