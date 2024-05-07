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

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

func TestGetAdapterRequest(t *testing.T) {
	tests := map[string]struct {
		req                *api_adapter_v1.GetPageRequest
		wantAdapterRequest *framework.Request[TestConfigA]
		wantReverseMapping *entityReverseIdMapping
		wantAdapterErr     *api_adapter_v1.Error
	}{
		"invalid_nil": {
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Request is nil.",
				Code:    1, // INVALID_PAGE_REQUEST_CONFIG
			},
		},
		"invalid_no_datasource_config": {
			req: &api_adapter_v1.GetPageRequest{
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
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Request contains no datasource config.",
				Code:    1, // INVALID_PAGE_REQUEST_CONFIG
			},
		},
		"invalid_no_entity_config": {
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
				PageSize: 100,
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Request contains no entity config.",
				Code:    1, // INVALID_PAGE_REQUEST_CONFIG
			},
		},
		"invalid_non_positive_page_size": {
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
				PageSize: 0,
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Request contains an invalid page size: 0. Must be greater than 0.",
				Code:    1, // INVALID_PAGE_REQUEST_CONFIG
			},
		},
		"invalid_datasource_config_no_id": {
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
					Ordered: true,
				},
				PageSize: 100,
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Datasource config contains no ID.",
				Code:    2, // INVALID_DATASOURCE_CONFIG
			},
		},
		"invalid_entity_config_config_not_json": {
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id:      "1f530a64-0565-49e6-8647-b88e908b7229",
					Config:  []byte(`invalid JSON`),
					Address: "http://example.com/",
					Auth: &api_adapter_v1.DatasourceAuthCredentials{
						AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
							HttpAuthorization: "Bearer mysecret",
						},
					},
				},
				Entity: &api_adapter_v1.EntityConfig{
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
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Config in datasource config could not parsed as JSON: invalid character 'i' looking for beginning of value.",
				Code:    2, // INVALID_DATASOURCE_CONFIG
			},
		},
		"invalid_entity_config_no_id": {
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
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Entity config contains no ID.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"all_fields_set": {
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
							List:       false,
							UniqueId:   true,
						},
					},
					Ordered: true,
				},
				PageSize: 100,
				Cursor:   "the cursor",
			},
			wantAdapterRequest: &framework.Request[TestConfigA]{
				Config: &TestConfigA{
					A: "a value",
					B: "b value",
				},
				Address: "http://example.com/",
				Auth: &framework.DatasourceAuthCredentials{
					HTTPAuthorization: "Bearer mysecret",
				},
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
							List:       false,
							UniqueId:   true,
						},
					},
				},
				Ordered:  true,
				PageSize: 100,
				Cursor:   "the cursor",
			},
			wantReverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       false,
						UniqueId:   true,
					},
				},
			},
		},
		"all_optional_fields_unset": {
			req: &api_adapter_v1.GetPageRequest{
				Datasource: &api_adapter_v1.DatasourceConfig{
					Id: "1f530a64-0565-49e6-8647-b88e908b7229",
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
			wantAdapterRequest: &framework.Request[TestConfigA]{
				Entity: framework.EntityConfig{
					ExternalId: "users",
					Attributes: []*framework.AttributeConfig{
						{
							ExternalId: "name",
							Type:       framework.AttributeTypeString,
						},
					},
				},
				PageSize: 100,
			},
			wantReverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotAdapterRequest, gotReverseMapping, gotAdapterErr := getAdapterRequest[TestConfigA](tc.req)
			AssertDeepEqual(t, tc.wantAdapterRequest, gotAdapterRequest)
			AssertDeepEqual(t, tc.wantReverseMapping, gotReverseMapping)
			AssertDeepEqual(t, tc.wantAdapterErr, gotAdapterErr)
		})
	}
}

func TestGetAdapterAuth(t *testing.T) {
	tests := map[string]struct {
		auth     *api_adapter_v1.DatasourceAuthCredentials
		wantAuth *framework.DatasourceAuthCredentials
	}{
		"basic": {
			auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_Basic_{
					Basic: &api_adapter_v1.DatasourceAuthCredentials_Basic{
						Username: "john",
						Password: "password123",
					},
				},
			},
			wantAuth: &framework.DatasourceAuthCredentials{
				Basic: &framework.BasicAuthCredentials{
					Username: "john",
					Password: "password123",
				},
			},
		},
		"http_authorization": {
			auth: &api_adapter_v1.DatasourceAuthCredentials{
				AuthMechanism: &api_adapter_v1.DatasourceAuthCredentials_HttpAuthorization{
					HttpAuthorization: "Bearer mysecret",
				},
			},
			wantAuth: &framework.DatasourceAuthCredentials{
				HTTPAuthorization: "Bearer mysecret",
			},
		},
		"nil": {
			auth:     nil,
			wantAuth: nil,
		},
		"all_fields_nil": {
			auth:     &api_adapter_v1.DatasourceAuthCredentials{},
			wantAuth: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotAuth := getAdapterAuth(tc.auth)
			AssertDeepEqual(t, tc.wantAuth, gotAuth)
		})
	}
}

func TestGetEntity(t *testing.T) {
	tests := map[string]struct {
		entity             *api_adapter_v1.EntityConfig
		wantAdapterEntity  *framework.EntityConfig
		wantReverseMapping *entityReverseIdMapping
		wantAdapterErr     *api_adapter_v1.Error
	}{
		"one_attribute": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
					},
				},
			},
			wantAdapterEntity: &framework.EntityConfig{
				ExternalId: "users",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "name",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
			},
			wantReverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
					},
				},
			},
		},
		"multiple_attributes": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
					},
					{
						Id:         "32b6245b-25a6-480b-9ca9-ac9bbe237a4e",
						ExternalId: "email",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "e660f6a5-eab8-4151-b56d-0ac259696d8a",
						ExternalId: "active",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					},
					{
						Id:         "1cdcb2a2-06e4-4637-bbec-55b5bc399579",
						ExternalId: "employeeNumber",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					},
					{
						Id:         "54e9ca5a-6948-4d9c-aaad-a673e6b1788b",
						ExternalId: "roles",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
			wantAdapterEntity: &framework.EntityConfig{
				ExternalId: "users",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "name",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
					{
						ExternalId: "email",
						Type:       framework.AttributeTypeString,
					},
					{
						ExternalId: "active",
						Type:       framework.AttributeTypeBool,
					},
					{
						ExternalId: "employeeNumber",
						Type:       framework.AttributeTypeInt64,
					},
					{
						ExternalId: "roles",
						Type:       framework.AttributeTypeString,
						List:       true,
					},
				},
			},
			wantReverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
					},
					"email": {
						Id:         "32b6245b-25a6-480b-9ca9-ac9bbe237a4e",
						ExternalId: "email",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					"active": {
						Id:         "e660f6a5-eab8-4151-b56d-0ac259696d8a",
						ExternalId: "active",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL,
					},
					"employeeNumber": {
						Id:         "1cdcb2a2-06e4-4637-bbec-55b5bc399579",
						ExternalId: "employeeNumber",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64,
					},
					"roles": {
						Id:         "54e9ca5a-6948-4d9c-aaad-a673e6b1788b",
						ExternalId: "roles",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						List:       true,
					},
				},
			},
		},
		"one_child_entity": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
								UniqueId:   true,
							},
						},
					},
				},
			},
			wantAdapterEntity: &framework.EntityConfig{
				ExternalId: "users",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "name",
						Type:       framework.AttributeTypeString,
					},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "entitlements",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "id",
								Type:       framework.AttributeTypeString,
								UniqueId:   true,
							},
						},
					},
				},
			},
			wantReverseMapping: &entityReverseIdMapping{
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
								UniqueId:   true,
							},
						},
					},
				},
			},
		},
		"multiple_child_entities": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
								UniqueId:   true,
							},
						},
					},
					{
						Id:         "a974da7c-48da-4262-8270-b83396abb563",
						ExternalId: "addresses",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
								UniqueId:   true,
							},
						},
					},
				},
			},
			wantAdapterEntity: &framework.EntityConfig{
				ExternalId: "users",
				Attributes: []*framework.AttributeConfig{
					{
						ExternalId: "name",
						Type:       framework.AttributeTypeString,
						UniqueId:   true,
					},
				},
				ChildEntities: []*framework.EntityConfig{
					{
						ExternalId: "entitlements",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "id",
								Type:       framework.AttributeTypeString,
								UniqueId:   true,
							},
						},
					},
					{
						ExternalId: "addresses",
						Attributes: []*framework.AttributeConfig{
							{
								ExternalId: "displayName",
								Type:       framework.AttributeTypeString,
								UniqueId:   true,
							},
						},
					},
				},
			},
			wantReverseMapping: &entityReverseIdMapping{
				Id: "00d58abb-0b80-4745-927a-af9b2fb612dd",
				Attributes: map[string]*api_adapter_v1.AttributeConfig{
					"name": {
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
						UniqueId:   true,
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
								UniqueId:   true,
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
								UniqueId:   true,
							},
						},
					},
				},
			},
		},
		"invalid_entity_no_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Entity config contains no ID.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_entity_no_external_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Entity config 00d58abb-0b80-4745-927a-af9b2fb612dd contains no external ID.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_entity_no_attributes": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains no attributes.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_no_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains no ID.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_no_external_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains no external ID: 12268f03-f99d-476f-91cc-5fe3404e1654.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_invalid_type": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       123,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains invalid type 123: 12268f03-f99d-476f-91cc-5fe3404e1654 (name).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_unspecified_type": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       0,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains unspecified type ATTRIBUTE_TYPE_UNSPECIFIED: 12268f03-f99d-476f-91cc-5fe3404e1654 (name).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_duplicate_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "email",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains duplicate ID: 12268f03-f99d-476f-91cc-5fe3404e1654 (email).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_attribute_duplicate_external_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
					{
						Id:         "32b6245b-25a6-480b-9ca9-ac9bbe237a4e",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Attribute in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains duplicate external ID: 32b6245b-25a6-480b-9ca9-ac9bbe237a4e (name).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_no_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Child entity in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains no ID.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_no_external_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Child entity in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains no external ID: 05182a15-2451-4551-80ef-606fd05c1cc2.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_duplicate_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "addresses",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Child entity in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains duplicate ID: 05182a15-2451-4551-80ef-606fd05c1cc2 (addresses).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_duplicate_external_id": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
					{
						Id:         "a974da7c-48da-4262-8270-b83396abb563",
						ExternalId: "entitlements",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "4b316f65f-2ce0-4577-8a77-74f52da4ad90",
								ExternalId: "displayName",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Child entity in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains duplicate external ID: a974da7c-48da-4262-8270-b83396abb563 (entitlements).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_same_external_id_as_attribute": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "name",
						Attributes: []*api_adapter_v1.AttributeConfig{
							{
								Id:         "41325064-39ac-4a67-994f-bdcc092642e4",
								ExternalId: "id",
								Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
							},
						},
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Child entity in entity config 00d58abb-0b80-4745-927a-af9b2fb612dd (users) contains same external ID as attribute: 05182a15-2451-4551-80ef-606fd05c1cc2 (name).",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
		"invalid_child_entity_no_attributes": {
			entity: &api_adapter_v1.EntityConfig{
				Id:         "00d58abb-0b80-4745-927a-af9b2fb612dd",
				ExternalId: "users",
				Attributes: []*api_adapter_v1.AttributeConfig{
					{
						Id:         "12268f03-f99d-476f-91cc-5fe3404e1654",
						ExternalId: "name",
						Type:       api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING,
					},
				},
				ChildEntities: []*api_adapter_v1.EntityConfig{
					{
						Id:         "05182a15-2451-4551-80ef-606fd05c1cc2",
						ExternalId: "entitlements",
					},
				},
			},
			wantAdapterErr: &api_adapter_v1.Error{
				Message: "Entity config 05182a15-2451-4551-80ef-606fd05c1cc2 (entitlements) contains no attributes.",
				Code:    4, // INVALID_ENTITY_CONFIG
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotAdapterEntity, gotReverseMapping, gotAdapterErr := getEntity(tc.entity)
			AssertDeepEqual(t, tc.wantAdapterEntity, gotAdapterEntity)
			AssertDeepEqual(t, tc.wantReverseMapping, gotReverseMapping)
			AssertDeepEqual(t, tc.wantAdapterErr, gotAdapterErr)
		})
	}
}
