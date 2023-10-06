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
	"testing"
)

func TestParseConfig(t *testing.T) {
	var wantConfig, gotConfig *TestConfigA
	var err error

	gotConfig, err = ParseConfig[TestConfigA](nil)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	wantConfig = nil
	AssertDeepEqual(t, wantConfig, gotConfig)

	gotConfig, err = ParseConfig[TestConfigA]([]byte(`{"a":"a value","b":"b value"}`))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	wantConfig = &TestConfigA{
		A: "a value",
		B: "b value",
	}
	AssertDeepEqual(t, wantConfig, gotConfig)

	_, err = ParseConfig[TestConfigA]([]byte(`invalid JSON`))
	if err == nil || err.Error() != "invalid character 'i' looking for beginning of value" {
		t.Errorf("Expected %v, got %v", "invalid character 'i' looking for beginning of value", err)
	}
}
