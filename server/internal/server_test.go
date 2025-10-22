// Copyright 2023 SGNL.ai, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/connector"
	"github.com/sgnl-ai/adapter-framework/pkg/logs"
	"github.com/sgnl-ai/adapter-framework/pkg/logs/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/codes"
	grpc_metadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
)

type MockAdapterA struct {
	Response framework.Response
}

func (a *MockAdapterA) GetPage(ctx context.Context, request *framework.Request[TestConfigA]) framework.Response {
	return a.Response
}

func NewAdapterA(resp framework.Response) framework.Adapter[TestConfigA] {
	return &MockAdapterA{resp}
}

type MockAdapterB struct {
	Response framework.Response
}

func (a *MockAdapterB) GetPage(ctx context.Context, request *framework.Request[TestConfigB]) framework.Response {
	return a.Response
}

func NewAdapterB(resp framework.Response) framework.Adapter[TestConfigB] {
	return &MockAdapterB{resp}
}

func TestServer_GetPage(t *testing.T) {
	validTokens := []string{"dGhpc2lzYXRlc3R0b2tlbg==", "dGhpc2lzYWxzb2F0ZXN0dG9rZW4="}

	tests := map[string]struct {
		req                  *api_adapter_v1.GetPageRequest
		tokens               []string
		adapterResponse      framework.Response
		wantResp             *api_adapter_v1.GetPageResponse
		wantError            error
		ctxWithConnectorInfo bool
	}{
		"success": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"name": "Alice",
						},
						{
							"name": "Bob",
						},
					},
					NextCursor: "next cursor",
				},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Success{
					Success: &api_adapter_v1.Page{
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Alice"}},
										},
									},
								},
							},
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Bob"}},
										},
									},
								},
							},
						},
						NextCursor: "next cursor",
					},
				},
			},
		},
		"success_with_connector_info": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					ConnectorInfo: &api_adapter_v1.ConnectorInfo{
						Id: "test-connector-id",
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"name": "Alice",
						},
						{
							"name": "Bob",
						},
					},
					NextCursor: "next cursor",
				},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Success{
					Success: &api_adapter_v1.Page{
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Alice"}},
										},
									},
								},
							},
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Bob"}},
										},
									},
								},
							},
						},
						NextCursor: "next cursor",
					},
				},
			},
		},
		"error_when_connector_with_context_fails": {
			ctxWithConnectorInfo: true,
			tokens:               []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					ConnectorInfo: &api_adapter_v1.ConnectorInfo{
						Id: "test-connector-id",
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Error creating connector context, context is already configured with the connector info, {   0 }.",
						Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL, // INVALID_DATASOURCE_CONFIG
					},
				},
			},
		},
		"success_explicit_type": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					Type: "Mock-1.0.1",
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"name": "Alice",
						},
						{
							"name": "Bob",
						},
					},
					NextCursor: "next cursor",
				},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Success{
					Success: &api_adapter_v1.Page{
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Alice"}},
										},
									},
								},
							},
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "Bob"}},
										},
									},
								},
							},
						},
						NextCursor: "next cursor",
					},
				},
			},
		},
		"error": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Error: &framework.Error{
					Message:    "Some error message.",
					Code:       api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
					RetryAfter: Ptr(23 * time.Second),
				},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message:    "Some error message.",
						Code:       api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
						RetryAfter: durationpb.New(23 * time.Second),
					},
				},
			},
		},
		"invalid_request": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
				},
				PageSize: 100,
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Datasource config contains no ID.",
						Code:    2, // INVALID_DATASOURCE_CONFIG
					},
				},
			},
		},
		"invalid_type": {
			tokens: []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					Type: "Invalid-1.0.0",
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{},
			},
			wantResp: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Unsupported datasource type provided: Invalid-1.0.0.",
						Code:    2, // INVALID_DATASOURCE_CONFIG
					},
				},
			},
		},
		"missing_auth_token": {
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					Type: "Invalid-1.0.0",
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{},
			},
			wantError: status.Errorf(codes.Unauthenticated, "invalid or missing token"),
		},
		"invalid_auth_token": {
			tokens: []string{"invalid"},
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`{"a":"a value","b":"b value"}`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
					Type: "Invalid-1.0.0",
				},
				Entity: &api_adapter_v1.EntityConfig{
					Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
					ExternalId: "users",
					Attributes: []*api_adapter_v1.AttributeConfig{
						{
							Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
							ExternalId: "name",
							Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			adapterResponse: framework.Response{
				Success: &framework.Page{},
			},
			wantError: status.Errorf(codes.Unauthenticated, "invalid or missing token"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			// Set the context with connector info to fail connector.WithContext() call
			// while processing the request.
			if tc.ctxWithConnectorInfo {
				ctx, _ = connector.WithContext(ctx, connector.ConnectorInfo{})
			}
			ctx = grpc_metadata.NewIncomingContext(ctx, grpc_metadata.MD{
				"token": tc.tokens,
			})

			s := &Server{
				Tokens:              validTokens,
				AdapterGetPageFuncs: make(map[string]AdapterGetPageFunc),
			}

			if err := RegisterAdapter(s, "Mock-1.0.1", NewAdapterA(tc.adapterResponse)); err != nil {
				t.Fatal(err)
			}

			if err := RegisterAdapter(s, "", NewAdapterB(tc.adapterResponse)); err != nil {
				t.Fatal(err)
			}

			gotResp, gotError := s.GetPage(ctx, tc.req)

			AssertDeepEqual(t, tc.wantResp, gotResp)
			AssertDeepEqual(t, tc.wantError, gotError)
		})
	}
}

type MockAdapterWithContext struct {
	Response    framework.Response
	CapturedCtx context.Context
}

func (a *MockAdapterWithContext) GetPage(ctx context.Context, request *framework.Request[TestConfigA]) framework.Response {
	a.CapturedCtx = ctx

	return a.Response
}

func TestServer_GetPage_WithLogger(t *testing.T) {
	validTokens := []string{"dGhpc2lzYXRlc3R0b2tlbg=="}

	mockAdapter := &MockAdapterWithContext{
		Response: framework.Response{
			Success: &framework.Page{
				Objects: []framework.Object{
					{"name": "Alice"},
				},
			},
		},
	}

	ctx := context.Background()
	ctx = grpc_metadata.NewIncomingContext(ctx, grpc_metadata.MD{
		"token": validTokens,
	})

	// Create an observable logger to capture log output.
	observedCore, observedLogs := observer.New(zapcore.InfoLevel)
	zapLogger := zap.New(observedCore)
	logger := zaplog.New(zapLogger)

	s := &Server{
		Tokens:              validTokens,
		AdapterGetPageFuncs: make(map[string]AdapterGetPageFunc),
		Logger:              logger,
	}

	if err := RegisterAdapter(s, "Mock-1.0.1", mockAdapter); err != nil {
		t.Fatal(err)
	}

	req := &api_adapter_v1.GetPageRequest{
		TenantId: "test-tenant-123",
		ClientId: "test-client-456",
		Datasource: &api_adapter_v1.DatasourceConfig{
			Id:      "datasource-789",
			Type:    "Mock-1.0.1",
			Config:  []byte(`{"a":"a value"}`),
			Address: "http://example.com/",
			Auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer mysecret",
				},
			},
		},
		Entity: &api_adapter_v1.EntityConfig{
			Id:         "entity-abc",
			ExternalId: "users",
			Attributes: []*api_adapter_v1.AttributeConfig{
				{
					Id:         "attr-123",
					ExternalId: "name",
					Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
				},
			},
			Ordered: true,
		},
		PageSize: 50,
		Cursor:   "test-cursor",
	}

	_, err := s.GetPage(ctx, req)
	if err != nil {
		t.Fatalf("GetPage returned error: %v", err)
	}

	// Verify logger was added to context
	if mockAdapter.CapturedCtx == nil {
		t.Fatal("Context was not passed to adapter")
	}

	retrievedLogger := logs.FromContext(mockAdapter.CapturedCtx)
	if retrievedLogger == nil {
		t.Fatal("Logger was not added to context")
	}

	// Test that the logger has the correct fields by logging a message and checking the output.
	retrievedLogger.Info("test message")

	if observedLogs.Len() != 1 {
		t.Fatalf("Expected 1 log entry, got %d", observedLogs.Len())
	}

	logEntry := observedLogs.All()[0]

	if logEntry.Message != "test message" {
		t.Errorf("Expected message 'test message', got %q", logEntry.Message)
	}

	// Verify the expected fields are present with correct values.
	expectedFields := map[string]any{
		"adapterRequestCursor":   "test-cursor",
		"adapterRequestPageSize": int64(50),
		"tenantId":               "test-tenant-123",
		"clientId":               "test-client-456",
		"datasourceAddress":      "http://example.com/",
		"datasourceId":           "datasource-789",
		"datasourceType":         "Mock-1.0.1",
		"entityId":               "entity-abc",
		"entityExternalId":       "users",
	}

	contextMap := logEntry.ContextMap()
	for fieldName, expectedValue := range expectedFields {
		actualValue, ok := contextMap[fieldName]
		if !ok {
			t.Errorf("Expected field %s not found in log output", fieldName)

			continue
		}

		if actualValue != expectedValue {
			t.Errorf("Field %s: expected %v, got %v", fieldName, expectedValue, actualValue)
		}
	}
}
