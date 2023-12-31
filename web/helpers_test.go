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
	"reflect"
	"testing"
	"time"
)

func AssertDeepEqual(t *testing.T, want, got any) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %#v, got %#v", want, got)
	}
}

func MustParseTime(t *testing.T, v string) time.Time {
	ts, err := time.Parse(time.RFC3339, v)
	if err != nil {
		t.Fatalf("Failed to parse %s as an RFC 3339 timestamp: %s", v, err)
	}
	return ts
}
