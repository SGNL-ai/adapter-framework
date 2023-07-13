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
	complexAttributeNameDelimiter string
	dateTimeFormats               []string
}

func defaultJSONOptions() *jsonOptions {
	return &jsonOptions{
		complexAttributeNameDelimiter: "", // Disabled.
		dateTimeFormats: []string{
			time.RFC3339, time.RFC3339Nano,
			time.RFC1123, time.RFC1123Z,
			time.RFC822, time.RFC822Z, time.RFC850,
			time.UnixDate, time.RubyDate, time.ANSIC,
			"2006-01-02T15:04:05.000Z0700",
			"2006-01-02 15:04:05", "2006-01-02",
			"2006/01/02", "01-02-2006", "01/02/2006",
			"01/02/06",
		},
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
// {
//   "attr1": {
//     "attr2": {
//       "attr3": "the value"
//     }
//   }
// }
//
// then the value returned for the attribute with external ID
// "attr1.attr2.attr3" is "the value".
//
// If empty (default), single-valued complex object parsing is disabled.
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
func WithDateTimeFormats(formats ...string) JSONOption {
	return &funcJSONOption{
		f: func(jo *jsonOptions) {
			jo.dateTimeFormats = formats
		},
	}
}
