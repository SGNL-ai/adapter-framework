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
		list, ok := value.([]any)
		if !ok {
			return nil, fmt.Errorf("list attribute %s cannot be parsed into a list value", attribute.ExternalId)
		}

		if len(list) == 0 {
			return list, nil
		}

		parsedList := make([]any, 0, len(list))

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

			parsedList = append(parsedList, parsedElement)
		}

		return parsedList, nil
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
		}
		return boolValue, nil

	case framework.AttributeTypeDateTime:
		v, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into a string date-time value", attribute.ExternalId)
		}
		if v == "" {
			return nil, nil
		}
		t, err := ParseDateTime(opts.dateTimeFormats, v)
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
		// Duration is assumed to be a number of seconds, as a float64,
		// possibly with a fractional part, e.g.: 3.5.
		v, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into a float64 duration value", attribute.ExternalId)
		}
		return time.Duration(v * float64(time.Second)), nil

	case framework.AttributeTypeInt64:
		// All numbers are unmarshalled into float64. Convert into int64.
		v, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("attribute %s cannot be parsed into an int64 value", attribute.ExternalId)
		}
		return v, nil

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

// ParseDateTime parses a timestamp against a set of predefined formats.
func ParseDateTime(formats []string, dateTimeStr string) (dateTime time.Time, err error) {
	for _, format := range formats {
		dateTime, err = time.Parse(format, dateTimeStr)
		if err == nil {
			return
		}
	}

	err = fmt.Errorf("failed to parse date-time value: %s", dateTimeStr)

	return
}
