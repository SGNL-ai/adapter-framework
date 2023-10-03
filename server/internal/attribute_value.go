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
	"fmt"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	nullValue = &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_NullValue{
		NullValue: &emptypb.Empty{},
	}}
)

func validateAttributeValue(attribute *api_adapter_v1.AttributeConfig, value any) (adapterErr *api_adapter_v1.Error) {
	if value == nil {
		return nil
	}

	valid := false
	switch value.(type) {
	case []any:
		value, _ := value.([]any)
		valid = len(value) == 0 && attribute.List
	case []bool:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL && attribute.List
	case []*bool:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL && attribute.List
	case []time.Time:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME && attribute.List
	case []*time.Time:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME && attribute.List
	case []framework.Duration:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DURATION && attribute.List
	case []*framework.Duration:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DURATION && attribute.List
	case []float64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DOUBLE && attribute.List
	case []*float64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DOUBLE && attribute.List
	case []int64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64 && attribute.List
	case []*int64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64 && attribute.List
	case []string:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING && attribute.List
	case []*string:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING && attribute.List
	case bool:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL && !attribute.List
	case *bool:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_BOOL && !attribute.List
	case time.Time:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME && !attribute.List
	case *time.Time:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DATE_TIME && !attribute.List
	case framework.Duration:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DURATION && !attribute.List
	case *framework.Duration:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DURATION && !attribute.List
	case float64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DOUBLE && !attribute.List
	case *float64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_DOUBLE && !attribute.List
	case int64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64 && !attribute.List
	case *int64:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_INT64 && !attribute.List
	case string:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING && !attribute.List
	case *string:
		valid = attribute.Type == api_adapter_v1.AttributeType_ATTRIBUTE_TYPE_STRING && !attribute.List
	}

	if !valid {
		return &api_adapter_v1.Error{
			Message: fmt.Sprintf("Adapter returned a value with invalid type %T for attribute %s (%s) with type %s (list=%t). This is always indicative of a bug within the Adapter implementation.",
				value, attribute.Id, attribute.ExternalId, attribute.Type, attribute.List),
			Code: api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}

	return nil
}

// getAttributeValues converts the singleton value or list of values for an
// attribute.
// Returns an error if the value's type is invalid.
func getAttributeValues(value any) (list []*api_adapter_v1.AttributeValue, adapterErr *api_adapter_v1.Error) {
	switch v := value.(type) {
	case []any:
		// validateAttributeValue only allows empty []any lists.
		// So if the value has already been validated, we know it must be empty and we return an empty list.
		// If v is not empty, getAttributeListValues will call getAttributeValue which will return an error.
		return getAttributeListValues(v)
	case []bool:
		return getAttributeListValues(v)
	case []*bool:
		return getAttributeListValues(v)
	case []time.Time:
		return getAttributeListValues(v)
	case []*time.Time:
		return getAttributeListValues(v)
	case []framework.Duration:
		return getAttributeListValues(v)
	case []*framework.Duration:
		return getAttributeListValues(v)
	case []float64:
		return getAttributeListValues(v)
	case []*float64:
		return getAttributeListValues(v)
	case []int64:
		return getAttributeListValues(v)
	case []*int64:
		return getAttributeListValues(v)
	case []string:
		return getAttributeListValues(v)
	case []*string:
		return getAttributeListValues(v)
	default: // Non-list attribute value.
		var singleValue *api_adapter_v1.AttributeValue

		singleValue, adapterErr = getAttributeValue(v)
		if adapterErr != nil {
			return
		}

		if singleValue != nullValue {
			list = []*api_adapter_v1.AttributeValue{singleValue}
		}

		return
	}
}

// getAttributeListValues converts a list of values for an attribute.
// Returns an error if the value is a list or if its type is invalid.
func getAttributeListValues[Element any](listValue []Element) (list []*api_adapter_v1.AttributeValue, adapterErr *api_adapter_v1.Error) {
	if listValue == nil {
		return nil, nil
	}

	list = make([]*api_adapter_v1.AttributeValue, len(listValue))
	for i, e := range listValue {
		list[i], adapterErr = getAttributeValue(e)
		if adapterErr != nil {
			return nil, adapterErr
		}
	}

	return
}

// getAttributeValue converts a single value for an attribute.
// Returns an error if the value is a list or if its type is invalid.
func getAttributeValue(value any) (*api_adapter_v1.AttributeValue, *api_adapter_v1.Error) {
	if value == nil {
		return nullValue, nil
	}

	switch v := value.(type) {
	case bool:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_BoolValue{
			BoolValue: v,
		}}, nil
	case *bool:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	case time.Time:
		_, timezoneOffset := v.Zone()
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_DatetimeValue{
			DatetimeValue: &api_adapter_v1.DateTime{
				Timestamp:      timestamppb.New(v),
				TimezoneOffset: int32(timezoneOffset),
			},
		}}, nil
	case *time.Time:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	case framework.Duration:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_DurationValue{
			DurationValue: &api_adapter_v1.Duration{Nanos: v.Nanos, Seconds: v.Seconds, Days: v.Days, Months: v.Months},
		}}, nil
	case *framework.Duration:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	case float64:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_DoubleValue{
			DoubleValue: v,
		}}, nil
	case *float64:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	case int64:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_Int64Value{
			Int64Value: v,
		}}, nil
	case *int64:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	case string:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_StringValue{
			StringValue: v,
		}}, nil
	case *string:
		if v == nil {
			return nullValue, nil
		}
		return getAttributeValue(*v)
	default:
		return nil, &api_adapter_v1.Error{
			Message: fmt.Sprintf("Adapter returned an attribute value with invalid type: %T. This is always indicative of a bug within the Adapter implementation.", value),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
}
