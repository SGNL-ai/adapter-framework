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

package web

import (
	"encoding/json"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
)

func TestConvertJSONAttributeValue(t *testing.T) {
	tests := map[string]struct {
		attribute *framework.AttributeConfig
		valueJSON string
		opts      *jsonOptions
		wantValue any
		wantError error
	}{
		"null": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeString,
			},
			valueJSON: `null`,
			wantValue: nil,
		},
		"bool": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
			},
			valueJSON: `true`,
			wantValue: true,
		},
		"bool_from_string_true": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
			},
			valueJSON: `"true"`,
			wantValue: true,
		},
		"bool_from_string_false": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
			},
			valueJSON: `"false"`,
			wantValue: false,
		},
		"bool_from_number_true": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
			},
			valueJSON: `1`,
			wantValue: true,
		},
		"bool_from_number_false": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
			},
			valueJSON: `0`,
			wantValue: false,
		},
		"bool_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeBool,
				List:       true,
			},
			valueJSON: `[true, 0, "true"]`,
			wantValue: []bool{true, false, true},
		},
		"datetime": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
			},
			valueJSON: `"2023-06-23T12:34:56-07:00"`,

			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}},
			wantValue: MustParseTime(t, "2023-06-23T12:34:56-07:00"),
		},
		"datetime_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
				List:       true,
			},
			valueJSON: `["2023-06-23T12:34:56-07:00", "2023-06-23T12:34:58+05:00"]`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}},
			wantValue: []time.Time{MustParseTime(t, "2023-06-23T12:34:56-07:00"), MustParseTime(t, "2023-06-23T12:34:58+05:00")},
		},
		"datetime_missing_tz": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
			},
			valueJSON: `"2023-06-23 12:34:56"`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{"2006-01-02 15:04:05", false}}},
			wantValue: MustParseTime(t, "2023-06-23T12:34:56Z"),
		},
		"date": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
			},
			valueJSON: `"2023-06-23"`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}},
			wantValue: MustParseTime(t, "2023-06-23T00:00:00Z"),
		},
		"date_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
				List:       true,
			},
			valueJSON: `["2023-06-23", "2023-06-23"]`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}},
			wantValue: []time.Time{MustParseTime(t, "2023-06-23T00:00:00Z"), MustParseTime(t, "2023-06-23T00:00:00Z")},
		},
		"date_with_neg_tz_offset": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
			},
			valueJSON: `"2023-06-23"`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}, localTimeZoneOffset: -10},
			wantValue: MustParseTime(t, "2023-06-23T00:00:00-10:00"),
		},
		"date_with_pos_tz_offset": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDateTime,
			},
			valueJSON: `"2023-06-23"`,
			opts:      &jsonOptions{dateTimeFormats: []DateTimeFormatWithTz{{time.RFC3339, true}, {"2006-01-02", false}}, localTimeZoneOffset: +4},
			wantValue: MustParseTime(t, "2023-06-23T00:00:00+04:00"),
		},
		"double": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDouble,
			},
			valueJSON: `123`,
			wantValue: float64(123.0),
		},
		"double_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDouble,
				List:       true,
			},
			valueJSON: `[12, 34, 56]`,
			wantValue: []float64{12, 34, 56},
		},
		"duration_from_string": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDuration,
			},
			valueJSON: `"123s"`,
			wantValue: 123 * time.Second,
		},
		"duration_from_number_of_seconds": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDuration,
			},
			valueJSON: `123`,
			wantValue: 123 * time.Second,
		},
		"duration_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeDuration,
				List:       true,
			},
			valueJSON: `[12, "34s", 56]`,
			wantValue: []time.Duration{12 * time.Second, 34 * time.Second, 56 * time.Second},
		},
		"int64": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeInt64,
			},
			valueJSON: `123`,
			wantValue: int64(123),
		},
		"int64_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeInt64,
				List:       true,
			},
			valueJSON: `[12, 34, 56]`,
			wantValue: []int64{12, 34, 56},
		},
		"string": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeString,
			},
			valueJSON: `"abc"`,
			wantValue: "abc",
		},
		"string_list": {
			attribute: &framework.AttributeConfig{
				ExternalId: "a",
				Type:       framework.AttributeTypeString,
				List:       true,
			},
			valueJSON: `["a", "b", "c"]`,
			wantValue: []string{"a", "b", "c"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var value any
			err := json.Unmarshal([]byte(tc.valueJSON), &value)
			if err != nil {
				t.Fatalf("Failed to unmarshal test input JSON value: %v", err)
			}

			gotValue, gotError := convertJSONAttributeValue(tc.attribute, value, tc.opts)
			AssertDeepEqual(t, tc.wantValue, gotValue)
			AssertDeepEqual(t, tc.wantError, gotError)
		})
	}
}
