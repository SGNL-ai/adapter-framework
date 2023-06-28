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
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestGetResponse(t *testing.T) {
	tests := map[string]struct {
		reverseMapping  *entityReverseIdMapping
		resp            *framework.Response
		wantRpcResponse *api_adapter_v1.GetPageResponse
	}{
		"error": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			resp: &framework.Response{
				Error: &framework.Error{
					Message:    "Some error message.",
					Code:       api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
					RetryAfter: Ptr(23 * time.Second),
				},
			},
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message:    "Some error message.",
						Code:       api_adapter_v1.ErrorCode_ERROR_CODE_DATASOURCE_FAILED,
						RetryAfter: durationpb.New(23 * time.Second),
					},
				},
			},
		},
		"success_no_objects": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
					"addresses": {
						Id: "a974da7c-48da-4262-8270-b83396abb563",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"displayName": {
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			resp: &framework.Response{
				Success: &framework.Page{
					Objects: nil,
				},
			},
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Success{
					Success: &api_adapter_v1.Page{
						Objects: nil,
					},
				},
			},
		},
		"success_multiple_objects": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
					"addresses": {
						Id: "a974da7c-48da-4262-8270-b83396abb563",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"displayName": {
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			resp: &framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						{
							"name": "Alice",
							"entitlements": []framework.Object{
								{
									"id": "id123",
								},
								{
									"id": "id456",
								},
							},
						},
						{
							"name": "Bob",
							"addresses": []framework.Object{
								{
									"displayName": "home",
								},
							},
						},
					},
					NextCursor: "nextCursor",
				},
			},
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
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
								ChildObjects: []*api_adapter_v1.EntityObjects{
									{
										EntityId: "05182a15-2451-4551-80ef-606fd05c1cc2",
										Objects: []*api_adapter_v1.Object{
											{
												Attributes: []*api_adapter_v1.Attribute{
													{
														Id: "41325064-39ac-4a67-994f-bdcc092642e4",
														Values: []*api_adapter_v1.AttributeValue{
															{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "id123"}},
														},
													},
												},
											},
											{
												Attributes: []*api_adapter_v1.Attribute{
													{
														Id: "41325064-39ac-4a67-994f-bdcc092642e4",
														Values: []*api_adapter_v1.AttributeValue{
															{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "id456"}},
														},
													},
												},
											},
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
								ChildObjects: []*api_adapter_v1.EntityObjects{
									{
										EntityId: "a974da7c-48da-4262-8270-b83396abb563",
										Objects: []*api_adapter_v1.Object{
											{
												Attributes: []*api_adapter_v1.Attribute{
													{
														Id: "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
														Values: []*api_adapter_v1.AttributeValue{
															{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "home"}},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						NextCursor: "nextCursor",
					},
				},
			},
		},
		"invalid_nil_response": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			resp: nil,
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Adapter returned nil response. This is always indicative of a bug within the Adapter implementation.",
						Code:    11, // ERROR_CODE_INTERNAL
					},
				},
			},
		},
		"invalid_empty_response": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			resp: &framework.Response{},
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Adapter returned empty response. This is always indicative of a bug within the Adapter implementation.",
						Code:    11, // ERROR_CODE_INTERNAL
					},
				},
			},
		},
		"invalid_objects": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			resp: &framework.Response{
				Success: &framework.Page{
					Objects: []framework.Object{
						nil, // Invalid.
					},
				},
			},
			wantRpcResponse: &api_adapter_v1.GetPageResponse{
				Response: &api_adapter_v1.GetPageResponse_Error{
					Error: &api_adapter_v1.Error{
						Message: "Adapter returned an object for entity 00d58abb-0b80-4745-927a-af9b2fb612dd which contains no non-null attributes. This is always indicative of a bug within the Adapter implementation.",
						Code:    11, // ERROR_CODE_INTERNAL
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotRpcResponse := getResponse(tc.reverseMapping, tc.resp)
			if gotRpcResponse.GetError() != nil {
				t.Logf("ERROR: %s", gotRpcResponse.GetError().Message)
			}
			AssertDeepEqual(t, tc.wantRpcResponse, gotRpcResponse)
		})
	}
}

func TestGetEntityObject(t *testing.T) {
	tests := map[string]struct {
		reverseMapping   *entityReverseIdMapping
		object           framework.Object
		wantEntityObject *api_adapter_v1.Object
		wantAdapterErr   *api_adapter_v1.Error
	}{
		"invalid_empty_object": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: nil,
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned an object for entity 00d58abb-0b80-4745-927a-af9b2fb612dd which contains no non-null attributes. This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"one_attribute_string": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: framework.Object{
				"name": "John Doe",
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
						},
					},
				},
			},
		},
		"one_attribute_string_list_multiple_values": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
			object: framework.Object{
				"name": []string{"John Doe", "John C. Doe"},
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John C. Doe"}},
						},
					},
				},
			},
		},
		"one_attribute_string_list_empty": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
			object: framework.Object{
				"name": []string{},
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id:     "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{},
					},
				},
			},
		},
		"one_attribute_string_list_containing_null": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
			object: framework.Object{
				"name": []*string{Ptr("John Doe"), nil},
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
							nullValue,
						},
					},
				},
			},
		},
		"null_attribute_skipped": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					"email": {
						Id:         "32b6245b-25a6-480b-9ca9-ac9bbe237a4e",
						ExternalId: "email",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					"roles": {
						Id:         "54e9ca5a-6948-4d9c-aaad-a673e6b1788b",
						ExternalId: "roles",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
			object: framework.Object{
				"name":  "John Doe",
				"email": (*string)(nil),  // null value; ignored
				"roles": ([]string)(nil), // null value; ignored
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
						},
					},
				},
			},
		},
		"invalid_no_non_null_attribute_value": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: framework.Object{
				"name": (*string)(nil),
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned an object for entity 00d58abb-0b80-4745-927a-af9b2fb612dd which contains no non-null attributes. This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"invalid_attribute_external_id": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: framework.Object{
				"email": "john@doe.org",
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned an object for entity 00d58abb-0b80-4745-927a-af9b2fb612dd which contains an attribute with an invalid external ID: email. This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"invalid_attribute_value_type": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: framework.Object{
				"name": int64(1234),
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned a value with invalid type int64 for attribute 12268f03-f99d-476f-91cc-5fe3404e1654 (name) with type ATTRIBUTE_TYPE_STRING (list=false). This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"one_child_entity": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			object: framework.Object{
				"name": "John Doe",
				"entitlements": []framework.Object{
					{
						"id": "id123",
					},
				},
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
						},
					},
				},
				ChildObjects: []*api_adapter_v1.EntityObjects{
					{
						EntityId: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "41325064-39ac-4a67-994f-bdcc092642e4",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "id123"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"one_child_entity_no_objects": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			object: framework.Object{
				"name":         "John Doe",
				"entitlements": []framework.Object{}, // Empty list of child objects; ignored
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
						},
					},
				},
			},
		},
		"multiple_child_entities": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
					"addresses": {
						Id: "a974da7c-48da-4262-8270-b83396abb563",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"displayName": {
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			object: framework.Object{
				"name": "John Doe",
				"entitlements": []framework.Object{
					{
						"id": "id123",
					},
				},
				"addresses": []framework.Object{
					{
						"displayName": "home",
					},
					{
						"displayName": "work",
					},
				},
			},
			wantEntityObject: &api_adapter_v1.Object{
				Attributes: []*api_adapter_v1.Attribute{
					{
						Id: "12268f03-f99d-476f-91cc-5fe3404e1654",
						Values: []*api_adapter_v1.AttributeValue{
							{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "John Doe"}},
						},
					},
				},
				ChildObjects: []*api_adapter_v1.EntityObjects{
					{
						EntityId: "a974da7c-48da-4262-8270-b83396abb563",
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "home"}},
										},
									},
								},
							},
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "work"}},
										},
									},
								},
							},
						},
					},
					{
						EntityId: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Objects: []*api_adapter_v1.Object{
							{
								Attributes: []*api_adapter_v1.Attribute{
									{
										Id: "41325064-39ac-4a67-994f-bdcc092642e4",
										Values: []*api_adapter_v1.AttributeValue{
											{Value: &api_adapter_v1.AttributeValue_StringValue{StringValue: "id123"}},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"invalid_child_object_external_id": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			object: framework.Object{
				"name": "John Doe",
				"entitlements": []framework.Object{
					{
						"id": "id123",
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned an object for entity 00d58abb-0b80-4745-927a-af9b2fb612dd which contains child objects with an invalid entity external ID: entitlements. This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"invalid_attribute_external_id_in_child_object": {
			reverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: map[string]*entityReverseIdMapping{
					"entitlements": {
						Id: "05182a15-2451-4551-80ef-606fd05c1cc2",
						Attributes: map[string]*api_adapter_v1.AttributeConfig{
							"id": {
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			object: framework.Object{
				"name": "John Doe",
				"entitlements": []framework.Object{
					{
						"displayName": "Some Entitlement",
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Adapter returned an object for entity 05182a15-2451-4551-80ef-606fd05c1cc2 which contains an attribute with an invalid external ID: displayName. This is always indicative of a bug within the Adapter implementation.",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotEntityObject, gotAdapterErr := getEntityObject(tc.reverseMapping, tc.object)
			AssertDeepEqual(t, tc.wantEntityObject, gotEntityObject)
			AssertDeepEqual(t, tc.wantAdapterErr, gotAdapterErr)
		})
	}
}
