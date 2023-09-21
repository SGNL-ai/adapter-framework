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
	"os"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	grpc_metadata "google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
)

type mockAdapter[Config any] struct {
	Response framework.Response
}

func (a *mockAdapter[Config]) GetPage(ctx context.Context, request *framework.Request[Config]) framework.Response {
	return a.Response
}

func TestServer_GetPage(t *testing.T) {
	path := "./TOKENS"
	if err := os.Setenv("AUTH_TOKENS_PATH", path); err != nil {
		t.Fatal(err)
	}

	token := []byte(`["dGhpc2lzYXRlc3R0b2tlbg==","dGhpc2lzYWxzb2F0ZXN0dG9rZW4="]`)
	if err := os.WriteFile(path, token, 0666); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		req             *api_adapter_v1.GetPageRequest
		adapterResponse framework.Response
		wantResp        *api_adapter_v1.GetPageResponse
	}{
		"success": {
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
		"success_explicit_type": {
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
		"unsupported_type": {
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
						Code:    14, // UNSUPPORTED_TYPE
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := grpc_metadata.NewIncomingContext(context.Background(), grpc_metadata.MD{
				"token": []string{"dGhpc2lzYXRlc3R0b2tlbg=="},
			})

			server := Server[TestConfig]{
				Adapters: map[string]framework.Adapter[TestConfig]{
					"":           &mockAdapter[TestConfig]{Response: tc.adapterResponse},
					"Mock-1.0.1": &mockAdapter[TestConfig]{Response: tc.adapterResponse},
				},
			}

			gotResp, _ := server.GetPage(ctx, tc.req)
			AssertDeepEqual(t, tc.wantResp, gotResp)
		})
	}
}
