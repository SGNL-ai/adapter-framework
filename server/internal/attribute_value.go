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

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// getAttributeValues converts the singleton value or list of values for an
// attribute.
// Returns an error if the value's type is invalid.
func getAttributeValues(value any) (list []*api_adapter_v1.AttributeValue, adapterErr *api_adapter_v1.Error) {
	switch v := value.(type) {
	case []bool:
		return getAttributeListValues(v)
	case []*bool:
		return getAttributeListValues(v)
	case []time.Time:
		return getAttributeListValues(v)
	case []*time.Time:
		return getAttributeListValues(v)
	case []time.Duration:
		return getAttributeListValues(v)
	case []*time.Duration:
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

		list = []*api_adapter_v1.AttributeValue{singleValue}

		return
	}
}

// getAttributeListValues converts a list of values for an attribute.
// Returns an error if the value is a list or if its type is invalid.
func getAttributeListValues[Element any](listValue []Element) (list []*api_adapter_v1.AttributeValue, adapterErr *api_adapter_v1.Error) {
	list = make([]*api_adapter_v1.AttributeValue, len(listValue))
	for i, e := range listValue {
		list[i], adapterErr = getAttributeValue(e)
		if adapterErr != nil {
			return
		}
	}

	return
}

// getAttributeValue converts a single value for an attribute.
// Returns an error if the value is a list or if its type is invalid.
func getAttributeValue(value any) (*api_adapter_v1.AttributeValue, *api_adapter_v1.Error) {
	if value == nil {
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_NullValue{
			NullValue: &emptypb.Empty{},
		}}, nil
	}

	switch v := value.(type) {
	case bool:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_BoolValue{
			BoolValue: v,
		}}, nil
	case *bool:
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
		return getAttributeValue(*v)
	case time.Duration:
		seconds := v / time.Second
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_DurationValue{
			DurationValue: &durationpb.Duration{
				Seconds: int64(seconds),
				Nanos:   int32(v - seconds*time.Second),
			},
		}}, nil
	case *time.Duration:
		return getAttributeValue(*v)
	case float64:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_DoubleValue{
			DoubleValue: v,
		}}, nil
	case *float64:
		return getAttributeValue(*v)
	case int64:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_Int64Value{
			Int64Value: v,
		}}, nil
	case *int64:
		return getAttributeValue(*v)
	case string:
		return &api_adapter_v1.AttributeValue{Value: &api_adapter_v1.AttributeValue_StringValue{
			StringValue: v,
		}}, nil
	case *string:
		return getAttributeValue(*v)
	default:
		return nil, &api_adapter_v1.Error{
			Message: fmt.Sprintf(api_adapter_v1.ErrorMsgAdapterInvalidAttributeValueTypeFmt, value),
			Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
		}
	}
}
