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
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
)

// convertJSONAttributeValue parses and converts the value of a JSON object
// field.
func convertJSONAttributeValue(attribute *framework.AttributeConfig, value any, opts *jsonOptions) (any, error) {
	if value == nil {
		return nil, nil
	}

	if attribute.List {
		switch attribute.Type {
		case framework.AttributeTypeBool:
			return convertJSONAttributeListValue[bool](attribute, value, opts)
		case framework.AttributeTypeDateTime:
			return convertJSONAttributeListValue[time.Time](attribute, value, opts)
		case framework.AttributeTypeDouble:
			return convertJSONAttributeListValue[float64](attribute, value, opts)
		case framework.AttributeTypeDuration:
			return convertJSONAttributeListValue[*framework.Duration](attribute, value, opts)
		case framework.AttributeTypeInt64:
			return convertJSONAttributeListValue[int64](attribute, value, opts)
		case framework.AttributeTypeString:
			return convertJSONAttributeListValue[string](attribute, value, opts)
		default:
			panic("invalid attribute type")
		}
	}

	switch attribute.Type {

	case framework.AttributeTypeBool:
		var boolValue bool
		switch v := value.(type) {
		case bool:
			boolValue = v
		case string:
			lowerCaseValue := strings.ToLower(v)
			switch lowerCaseValue {
			case "false", "0":
				boolValue = false
			case "true", "1":
				boolValue = true
			default:
				return nil, fmt.Errorf("attribute %s cannot be parsed into a bool value", attribute.ExternalId)
			}
		case float64:
			intValue := int(v)
			switch intValue {
			case 0:
				boolValue = false
			case 1:
				boolValue = true
			default:
				return nil, fmt.Errorf("attribute %s cannot be parsed into a bool value", attribute.ExternalId)
			}
		default:
			return nil, fmt.Errorf("attribute %s cannot be parsed into a bool value", attribute.ExternalId)
		}
		return boolValue, nil

	case framework.AttributeTypeDateTime:
		var dateTimeStr string
		switch v := value.(type) {
		case string:
			dateTimeStr = v
		case float64:
			if v > float64(math.MaxInt64) || v < float64(math.MinInt64) {
				return nil, fmt.Errorf("attribute %s cannot be parsed into a string date-time value because the float64 value is out of the int64 range", attribute.ExternalId)
			}
			dateTimeStr = fmt.Sprintf("%d", int(v))
		default:
			return nil, fmt.Errorf("attribute %s cannot be parsed into a string date-time value because its type is %T, not string or float64", attribute.ExternalId, v)
		}

		if dateTimeStr == "" {
			return nil, nil
		}
		t, err := ParseDateTime(opts.dateTimeFormats, opts.localTimeZoneOffset, dateTimeStr)
		if err != nil {
			return nil, fmt.Errorf("attribute %s cannot be parsed into a date-time value: %w", attribute.ExternalId, err)
		}
		return t, nil

	case framework.AttributeTypeDouble:
		v, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into a float64 value", attribute.ExternalId)
		}
		return v, nil

	case framework.AttributeTypeDuration:
		switch v := value.(type) {
		case string:
			duration, err := framework.ParseISO8601Duration(v)
			if err != nil {
				return nil, fmt.Errorf("attribute %s cannot be parsed into a duration value: %w", attribute.ExternalId, err)
			}
			return duration, nil
		default:
			return nil, fmt.Errorf("attribute %s cannot be parsed into a duration value", attribute.ExternalId)
		}

	case framework.AttributeTypeInt64:
		// All numbers are unmarshalled into float64. Convert into int64.
		v, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into an int64 value", attribute.ExternalId)
		}
		return int64(v), nil

	case framework.AttributeTypeString:
		v, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into a string value", attribute.ExternalId)
		}
		return v, nil

	default:
		panic("invalid attribute type")
	}
}

func convertJSONAttributeListValue[Element any](attribute *framework.AttributeConfig, value any, opts *jsonOptions) (any, error) {
	list, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("list attribute %s cannot be parsed into a list value", attribute.ExternalId)
	}

	if len(list) == 0 {
		return list, nil
	}

	parsedList := make([]Element, 0, len(list))

	elementAttribute := *attribute
	elementAttribute.List = false

	for _, element := range list {
		parsedElement, err := convertJSONAttributeValue(&elementAttribute, element, opts)

		if err != nil {
			return nil, err
		}

		// Do not return null attribute values.
		if parsedElement == nil {
			continue
		}

		parsedList = append(parsedList, parsedElement.(Element))
	}

	return parsedList, nil
}

// ParseDateTime parses a timestamp against a set of predefined formats.
func ParseDateTime(dateTimeFormats []DateTimeFormatWithTimeZone, localTimeZoneOffset int, dateTimeStr string) (dateTime time.Time, err error) {
	for _, format := range dateTimeFormats {
		if format.Format == "SGNLUnixSec" {
			var unixTimestamp int64
			unixTimestamp, err = strconv.ParseInt(dateTimeStr, 10, 64)
			if err == nil {
				dateTime = time.Unix(unixTimestamp, 0)
			}
		} else {
			dateTime, err = time.Parse(format.Format, dateTimeStr)
		}

		if err == nil {
			if !format.HasTimeZone {
				var loc *time.Location

				if localTimeZoneOffset == 0 {
					loc = time.UTC
				} else {
					loc = time.FixedZone("", localTimeZoneOffset)
				}

				dateTime = time.Date(
					dateTime.Year(),
					dateTime.Month(),
					dateTime.Day(),
					dateTime.Hour(),
					dateTime.Minute(),
					dateTime.Second(),
					dateTime.Nanosecond(),
					loc,
				)
			}

			return
		}
	}

	err = fmt.Errorf("failed to parse date-time value: %s", dateTimeStr)

	return
}
