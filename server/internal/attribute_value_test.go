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
	"encoding/json"
	"testing"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestGetAttributeListValues(t *testing.T) {
	testGetAttributeListValues[string](t, nil, "", nil)
	testGetAttributeListValues(t, []string{}, "", nil)
	testGetAttributeListValues(t, []string{"abcd"}, `[{"stringValue":"abcd"}]`, nil)
	testGetAttributeListValues(t, []string{"a", "b", "c"}, `[{"stringValue":"a"},{"stringValue":"b"},{"stringValue":"c"}]`, nil)
	testGetAttributeListValues(t, []*string{nil}, `[{"nullValue":{}}]`, nil)
	testGetAttributeListValues(t, []*string{Ptr("a"), nil, Ptr("c")}, `[{"stringValue":"a"},{"nullValue":{}},{"stringValue":"c"}]`, nil)
	testGetAttributeListValues(t, []int64{}, "", nil)
	testGetAttributeListValues(t, []int64{123}, `[{"int64Value":"123"}]`, nil)
	testGetAttributeListValues(t, []int64{12, 34, 56}, `[{"int64Value":"12"},{"int64Value":"34"},{"int64Value":"56"}]`, nil)
	testGetAttributeListValues(t, []*int64{nil}, `[{"nullValue":{}}]`, nil)
	testGetAttributeListValues(t, []*int64{Ptr(int64(12)), nil, Ptr(int64(56))}, `[{"int64Value":"12"},{"nullValue":{}},{"int64Value":"56"}]`, nil)

	testGetAttributeListValues(t, []int32{12, 34, 56}, "",
		&api_adapter_v1.Error{
			Message: "Adapter returned an attribute value with invalid type: int32",
			Code:    11, // ERROR_CODE_INTERNAL
		})
	testGetAttributeListValues(t, []**int64{Ptr(Ptr(int64(12))), nil, Ptr(Ptr(int64(56)))}, "",
		&api_adapter_v1.Error{
			Message: "Adapter returned an attribute value with invalid type: **int64",
			Code:    11, // ERROR_CODE_INTERNAL
		})
}

func testGetAttributeListValues[Element any](t *testing.T, listValue []Element, wantListAsJSON string, wantError *api_adapter_v1.Error) {
	t.Helper()

	var wantListOfMaps []map[string]any
	if wantListAsJSON != "" {
		err := json.Unmarshal([]byte(wantListAsJSON), &wantListOfMaps)
		if err != nil {
			t.Fatalf("Failed to unmarshal expected list of Protocol Buffer message: %v", err)
		}
	}

	gotList, gotError := getAttributeListValues(listValue)

	var gotListOfMaps []map[string]any
	if gotList != nil {
		gotListOfMaps = make([]map[string]any, 0, len(gotList))
		for _, gotAttributeValue := range gotList {
			var gotAttributeValueAsMap map[string]any
			if gotAttributeValue != nil {
				gotAttributeValueJSON, err := protojson.MarshalOptions{}.Marshal(gotAttributeValue)
				if err != nil {
					t.Fatalf("Failed to marshal Protocol Buffer message: %v", err)
				}

				gotAttributeValueAsMap = make(map[string]any)
				err = json.Unmarshal(gotAttributeValueJSON, &gotAttributeValueAsMap)
				if err != nil {
					t.Fatalf("Failed to unmarshal Protocol Buffer message: %v", err)
				}
			}
			gotListOfMaps = append(gotListOfMaps, gotAttributeValueAsMap)
		}
	}

	AssertDeepEqual(t, wantError, gotError)
	AssertDeepEqual(t, wantListOfMaps, gotListOfMaps)
}

func TestGetAttributeValue(t *testing.T) {
	timeValue, _ := time.Parse(time.RFC3339, "2023-06-23T12:34:56-07:00")

	tests := map[string]struct {
		value                  any
		wantAttributeValueJSON *string
		wantError              *api_adapter_v1.Error
	}{
		"null": {
			value:                  nil,
			wantAttributeValueJSON: Ptr(`{"nullValue":{}}`),
		},
		"bool": {
			value:                  true,
			wantAttributeValueJSON: Ptr(`{"boolValue":true}`),
		},
		"bool_pointer": {
			value:                  Ptr(true),
			wantAttributeValueJSON: Ptr(`{"boolValue":true}`),
		},
		"time": {
			value:                  timeValue,
			wantAttributeValueJSON: Ptr(`{"datetimeValue":{"timestamp":"2023-06-23T19:34:56Z", "timezoneOffset":-25200}}`),
		},
		"time_pointer": {
			value:                  Ptr(timeValue),
			wantAttributeValueJSON: Ptr(`{"datetimeValue":{"timestamp":"2023-06-23T19:34:56Z", "timezoneOffset":-25200}}`),
		},
		"duration": {
			value:                  12345 * time.Millisecond,
			wantAttributeValueJSON: Ptr(`{"durationValue":"12.345s"}`),
		},
		"duration_pointer": {
			value:                  Ptr(12345 * time.Millisecond),
			wantAttributeValueJSON: Ptr(`{"durationValue":"12.345s"}`),
		},
		"double": {
			value:                  float64(123.45),
			wantAttributeValueJSON: Ptr(`{"doubleValue":123.45}`),
		},
		"double_pointer": {
			value:                  Ptr(float64(123.45)),
			wantAttributeValueJSON: Ptr(`{"doubleValue":123.45}`),
		},
		"int64": {
			value:                  int64(1234),
			wantAttributeValueJSON: Ptr(`{"int64Value":"1234"}`),
		},
		"int64_pointer": {
			value:                  Ptr(int64(1234)),
			wantAttributeValueJSON: Ptr(`{"int64Value":"1234"}`),
		},
		"string": {
			value:                  "abcd",
			wantAttributeValueJSON: Ptr(`{"stringValue":"abcd"}`),
		},
		"string_pointer": {
			value:                  Ptr("abcd"),
			wantAttributeValueJSON: Ptr(`{"stringValue":"abcd"}`),
		},
		"invalid_int32": {
			value: int32(1234),
			wantError: &api_adapter_v1.Error{
				Message: "Adapter returned an attribute value with invalid type: int32",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
		"invalid_int64_pointer_pointer": {
			value: Ptr(Ptr(int64(1234))),
			wantError: &api_adapter_v1.Error{
				Message: "Adapter returned an attribute value with invalid type: **int64",
				Code:    11, // ERROR_CODE_INTERNAL
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Unmarshal the expected protobuf message marshaled into JSON so
			// that it can be compared using reflect.DeepEqual with the
			// protobuf message returned by getAttributeValue.
			// Using JSON in test cases is more readable than instantiating of
			// protobuf structs.
			var wantAttributeValueAsMap map[string]any
			if tc.wantAttributeValueJSON != nil {
				wantAttributeValueAsMap = make(map[string]any)
				err := json.Unmarshal([]byte(*tc.wantAttributeValueJSON), &wantAttributeValueAsMap)
				if err != nil {
					t.Fatalf("Failed to unmarshal expected Protocol Buffer message: %v", err)
				}
			}

			gotAttributeValue, gotError := getAttributeValue(tc.value)

			var gotAttributeValueAsMap map[string]any
			if gotAttributeValue != nil {
				gotAttributeValueJSON, err := protojson.MarshalOptions{}.Marshal(gotAttributeValue)
				if err != nil {
					t.Fatalf("Failed to marshal Protocol Buffer message: %v", err)
				}

				gotAttributeValueAsMap = make(map[string]any)
				err = json.Unmarshal(gotAttributeValueJSON, &gotAttributeValueAsMap)
				if err != nil {
					t.Fatalf("Failed to unmarshal Protocol Buffer message: %v", err)
				}
			}

			AssertDeepEqual(t, tc.wantError, gotError)
			AssertDeepEqual(t, wantAttributeValueAsMap, gotAttributeValueAsMap)
		})
	}
}
