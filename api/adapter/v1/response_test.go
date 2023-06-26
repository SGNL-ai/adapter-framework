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

package v1

import (
	"reflect"
	"testing"
)

func TestNewGetPageResponseSuccess(t *testing.T) {
	page := &Page{}
	wantResponse := &GetPageResponse{
		Response: &GetPageResponse_Success{
			Success: page,
		},
	}

	gotResponse := NewGetPageResponseSuccess(page)

	if !reflect.DeepEqual(wantResponse, gotResponse) {
		t.Errorf("Expected %#v, got %#v", wantResponse, gotResponse)
	}
}

func TestNewGetPageResponseError(t *testing.T) {
	err := &Error{}
	wantResponse := &GetPageResponse{
		Response: &GetPageResponse_Error{
			Error: err,
		},
	}

	gotResponse := NewGetPageResponseError(err)

	if !reflect.DeepEqual(wantResponse, gotResponse) {
		t.Errorf("Expected %#v, got %#v", wantResponse, gotResponse)
	}
}
