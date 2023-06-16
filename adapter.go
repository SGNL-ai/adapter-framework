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
package framework

import (
	"context"
	"time"

	v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// Adapter is the high-level interface implemented by adapters.
type Adapter[Config any] interface {
	// GetPage returns a page of objects from the requested datasource for the
	// requested entity.
	GetPage(ctx context.Context, request Request[Config]) Response
}

// Request is a request for a page of objects from a datasource for an entity.
//
// The Config type parameter must be a struct type into which the configuration
// JSON object can be unmarshaled into.
type Request[Config any] struct {
	// Config is configuration for the datasource.
	// Optional.
	Config *Config `json:"config,omitempty"`

	// Address is the address of the datasource.
	// Optional.
	Address string `json:"address,omitempty"`

	// Auth contains the credentials to use to authenticate with the
	// datasource.
	// Optional.
	Auth *DatasourceAuthCredentials `json:"auth,omitempty"`

	// EntityConfig is the configuration of the entity to get data from.
	EntityConfig EntityConfig `json:"entityConfig"`

	// Ordered indicates whether the objects are ordered by ID, i.e. whether
	// the response must contain objects ordered by monotonically increasing
	// IDs for the entity.
	// If true and the adapter cannot return objects ordered by ID, the adapter
	// must return error code ErrorCode_ERROR_CODE_INVALID_ENTITY_CONFIG.
	Ordered bool `json:"ordered,omitempty"`

	// PageSize is the maximum number of objects to return from the entity.
	PageSize int64 `json:"pageSize"`

	// Cursor identifies the first object of the page to return, as returned by
	// the last call to GetPage for the entity.
	// Optional. If not set, return the first page for this entity.
	Cursor string `json:"cursor,omitempty"`
}

// DatasourceAuthCredentials contains the credentials to authenticate with a
// datasource.
// Exactly one field is non-nil.
type DatasourceAuthCredentials struct {
	// Basic contains the credentials for basic username/password
	// authentication.
	Basic *BasicAuthCredentials `json:"basic,omitempty"`

	// HTTPAuthorization contains the credentials to be sent in an HTTP Authorization header.
	// Prefixed with the scheme, e.g. "Bearer ".
	HTTPAuthorization string `json:"httpAuthorization,omitempty"`
}

// BasicAuthCredentials contains credentials for basic username/password
// authentication with a datasource.
type BasicAuthCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// EntityConfig is the configuration of an entity to get data from.
type EntityConfig struct {
	// Id is the identifier of the entity within the datasource.
	Id string `json:"id"`

	// Attributes is the configuration of the attributes to return for the
	// entity.
	// Contains at least the entity's unique ID attribute.
	Attributes []AttributeConfig `json:"attributes"`

	// ChildEntities is the configuration of the entities that are children of
	// the entity to return together with the entity.
	// Optional.
	ChildEntities []EntityConfig `json:"childEntities,omitempty"`
}

// AttributeConfig is the configuration of an attribute to return.
type AttributeConfig struct {
	// Id is the adapter-specific name of the attribute in the entity.
	Id string `json:"id"`

	// Type is the type of the attribute's values.
	Type AttributeType `json:"type"`

	// List indicates whether the attribute contains a list of values vs. a
	// single value.
	List bool `json:"list,omitempty"`
}

// AttributeType is the type of the values for an attribute.
type AttributeType int32

const (
	// Boolean.
	AttributeTypeBool AttributeType = 1
	// Date or timestamp.
	AttributeTypeDateTime AttributeType = 2
	// Double-precision float.
	AttributeTypeDouble AttributeType = 3
	// Duration.
	AttributeTypeDuration AttributeType = 4
	// Signed integer.
	AttributeTypeInt64 AttributeType = 5
	// String.
	AttributeTypeString AttributeType = 6
)

// Response is the response to a GetPage request.
// Exactly one field must be non-nil.
type Response struct {
	Success *Page  `json:"success,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

// Page contains the objects requested by a GetPage request.
type Page struct {
	// Objects is the set of objects in the page returned by the datasource for
	// the requested entity.
	// Optional.
	//
	// Each object must be either:
	//  - a map[string]any which keys are attribute names,
	//  - or a pointer to an instance of a struct type.
	//
	// If it is an instance of a struct type, each field of that struct type
	// must have a tag containing either:
	//  - the "attr" tag key, which associated value is an attribute name,
	//  - or the "child" tag key, which associated value is a child entity ID.
	//
	// The type of a field tagged with "attr" must match the type requested for
	// that attribute in its AttributeConfig:
	//  - AttributeTypeBool -> bool
	//  - AttributeTypeDateTime -> time.Time
	//  - AttributeTypeDouble -> float64
	//  - AttributeTypeDuration -> time.Duration
	//  - AttributeTypeInt64 -> int64
	//  - AttributeTypeString -> string
	// The field type may be a pointer, if the attribute may have the nil
	// value.
	//
	// The type of a field tagged with "child" must be a list of map[string]any
	// maps or a list of struct values.
	//
	// For example:
	//
	// struct Example struct {
	//     Name      string    `attr:"name"`
	//     Email     *string   `attr:"email"`
	//     Addresses []Address `child:"address"`
	// }
	// struct Address struct {
	//     PostalCode *string `attr:"postalCode"`
	//     Region     *string `attr:"region"`
	// }
	//
	// Attributes that were not requested are ignored.
	Objects []any `json:"objects,omitempty"`

	// NextCursor the cursor that identifies the first object of the next page.
	// Optional. If not set, this page is the last page for this entity.
	NextCursor string `json:"nextCursor,omitempty"`
}

// Error contains the details of an error that occurred while executing a
// GetPage request.
type Error struct {
	// Message is the error message.
	// By convention, should start with an upper-case letter and end with a period.
	// Optional.
	Message string `json:"message,omitempty"`

	// Code is the error code indicating the cause of the error.
	Code v1.ErrorCode `json:"code"`

	// RetryAfter is the recommended minimal duration after which this request
	// may be retried.
	// Optional.
	RetryAfter *time.Duration `json:"retryAfter,omitempty"`
}
