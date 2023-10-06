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

package server

import (
	"reflect"
	"testing"
)

// TestConfigA is an example Config type which can be used as the Config type
// parameter of an Adapter for testing.
type TestConfigA struct {
	A string `json:"a"`
	B string `json:"b"`
}

// TestConfigB is an example Config type which can be used as the Config type
// parameter of an Adapter for testing.
type TestConfigB struct {
	C string `json:"c"`
	D string `json:"d"`
}

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// AssertDeepEqual asserts whether want and got are equal using
// reflect.DeepEqual.
func AssertDeepEqual(t *testing.T, want, got any) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %#v, got %#v", want, got)
	}
}
