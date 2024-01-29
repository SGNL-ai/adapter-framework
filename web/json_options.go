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

import "time"

// jsonOptions configures JSON object parsing. The fields are set by the
// JSONOption values passed to ConvertJSONObjectList.
type jsonOptions struct {
	// complexAttributeNameDelimiter is the delimiter to use to separate
	// hierarchical attribute names in attribute external IDs.
	// That feature is disabled if "".
	complexAttributeNameDelimiter string

	// dateTimeFormats is the set of datetime formats that are supported for
	// parsing the values of attributes of type datetime.
	dateTimeFormats []DateTimeFormatWithTimeZone

	// enableJSONPath indicates whether JSONPath is supported as a syntax for
	// attribute external IDs. If true, any external ID starting with '$' is
	// considered to be a JSONPath.
	enableJSONPath bool

	// localTimeZoneOffset is the default local timezone offset, as a number
	// of seconds east of UTC, to be used for parsing date-time attributes
	// lacking any timezone info.
	// Defaults to 0, i.e the UTC timezone.
	localTimeZoneOffset int
}

// DateTimeFormatWithTimeZone represents a valid date time format to try parsing
// date-time attribute values from strings.
type DateTimeFormatWithTimeZone struct {
	// Format must be a valid time format accepted by time.Parse.
	Format string

	// HasTimeZone indicates whether the above time format supports specifying a time zone.
	// If it does, this should be set to true.
	// If it does not, this can be set to false to use the specified localTimeZoneOffset
	// as a time zone in the resulting date time. If this value is false and
	// localTimeZoneOffset is not set, the resulting date time will be set to UTC.
	HasTimeZone bool
}

func defaultJSONOptions() *jsonOptions {
	return &jsonOptions{
		complexAttributeNameDelimiter: "", // Disabled.
		dateTimeFormats: []DateTimeFormatWithTimeZone{
			{time.RFC3339, true},
			{time.RFC3339Nano, true},
			{time.RFC1123Z, true},
			{time.RFC822, true},
			{time.RFC822Z, true},
			{time.RFC850, true},
			{time.UnixDate, true},
			{time.RubyDate, true},
			{"2006-01-02T15:04:05.000Z0700", true},
			{"2006-01-02 15:04:05", false},
			{"1/2/2006 3:04:05 PM", false},
			{time.ANSIC, false},
			{"2006-01-02", false},
			{"2006/01/02", false},
			{"01-02-2006", false},
			{"01/02/2006", false},
			{"01/02/06", false},
			{"SGNLUnixSec", false}, // Unix timestamp representing seconds since 1970-01-01 00:00:00 UTC.
		},
		enableJSONPath:      false, // Disabled.
		localTimeZoneOffset: 0,
	}
}

// JSONOption configures how JSON objects are parsed.
type JSONOption interface {
	apply(*jsonOptions)
}

type funcJSONOption struct {
	f func(*jsonOptions)
}

func (o *funcJSONOption) apply(opts *jsonOptions) {
	o.f(opts)
}

// WithComplexAttributeNameDelimiter sets the delimiter between nested objects
// names in attribute external IDs.
//
// If non-empty, and an attribute's external ID includes this delimiter,
// the attribute is parsed as an attribute of a single-valued complex
// object.
//
// For instance, if set to ".", and the JSON object to parse is:
//
//	{
//	  "attr1": {
//	    "attr2": {
//	      "attr3": "the value"
//	    }
//	  }
//	}
//
// then the value returned for the attribute with external ID
// "attr1.attr2.attr3" is "the value".
//
// If empty (default), single-valued complex object parsing is disabled.
//
// Deprecated: This option is replaced with WithJSONPathAttributeNames.
func WithComplexAttributeNameDelimiter(delimiter string) JSONOption {
	return &funcJSONOption{
		f: func(jo *jsonOptions) {
			jo.complexAttributeNameDelimiter = delimiter
		},
	}
}

// WithDateTimeFormats sets the time formats to use to try parsing date-time
// attribute values from strings.
// The formats must be ordered by decreasing likelihood of matching.
// Each format must be a valid time format accepted by time.Parse.
func WithDateTimeFormats(formats ...DateTimeFormatWithTimeZone) JSONOption {
	return &funcJSONOption{
		f: func(jo *jsonOptions) {
			jo.dateTimeFormats = formats
		},
	}
}

// WithJSONPathAttributeNames enables attribute external IDs specified as
// JSONPath to match attributes in nested objects.
//
// If enabled, and an attribute's external ID starts with '$', it is used as a
// JSONPath to match an attribute of a single- or multi-valued complex object.
//
// If the attribute is configured with a list type:
// if one or more values matched, then the resulting value is the list of
// matched values,
// otherwise (no value matched), then the resulting value is null.
//
// If the attribute is configured with a non-list type:
// if exactly one value matched, then the resulting value is that value,
// otherwise, if no value matched then the resulting value is null,
// otherwise (the JSONPath matched more than one value), an error is returned.
func WithJSONPathAttributeNames() JSONOption {
	return &funcJSONOption{
		f: func(jo *jsonOptions) {
			jo.enableJSONPath = true
		},
	}
}

// WithLocalTimeZoneOffset sets the local time zone offset to use as a default
// when parsing date-time attribute values from strings for formats lacking
// support for specifying a time zone.
func WithLocalTimeZoneOffset(offset int) JSONOption {
	return &funcJSONOption{
		f: func(jo *jsonOptions) {
			jo.localTimeZoneOffset = offset
		},
	}
}
