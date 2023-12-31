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
	"testing"
	"time"
)

func TestWithComplexAttributeNameDelimiter(t *testing.T) {
	var opts jsonOptions

	opt := WithComplexAttributeNameDelimiter("->")

	opt.apply(&opts)

	AssertDeepEqual(t, "->", opts.complexAttributeNameDelimiter)
}

func TestWithDateTimeFormats(t *testing.T) {
	var opts jsonOptions

	dateTimeFormats := []DateTimeFormatWithTimeZone{
		{time.RFC3339, true},
		{time.RFC3339Nano, true},
	}

	opt := WithDateTimeFormats(dateTimeFormats...)

	opt.apply(&opts)

	AssertDeepEqual(t, dateTimeFormats, opts.dateTimeFormats)
}

func TestWithJSONPathAttributeNames(t *testing.T) {
	var opts jsonOptions

	opt := WithJSONPathAttributeNames()

	opt.apply(&opts)

	AssertDeepEqual(t, true, opts.enableJSONPath)
}

func TestWithLocalTimeZoneOffset(t *testing.T) {
	var opts jsonOptions

	offset := -7 * 60 * 60

	opt := WithLocalTimeZoneOffset(offset)

	opt.apply(&opts)

	AssertDeepEqual(t, offset, opts.localTimeZoneOffset)
}
