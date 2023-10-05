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
	"os"
	"testing"
	"time"

	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server/internal"
)

func TestNewWithAuthTokensPath(t *testing.T) {
	validTokensPath := "./TOKENS_0"

	tokens := []byte(`["dGhpc2lzYXRlc3R0b2tlbg==","dGhpc2lzYWxzb2F0ZXN0dG9rZW4="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}

	invalidTokensPath := "./TOKENS_INVALID_0"

	invalidTokens := []byte(`invalidtokenformat`)
	if err := os.WriteFile(invalidTokensPath, invalidTokens, 0666); err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		inputAuthTokensPath string
		inputStopChan       <-chan struct{}
		wantAdapterServer   api_adapter_v1.AdapterServer
	}{
		"simple": {
			inputAuthTokensPath: validTokensPath,
			inputStopChan:       nil,
			wantAdapterServer: &internal.Server{
				Tokens:              []string{"dGhpc2lzYXRlc3R0b2tlbg==", "dGhpc2lzYWxzb2F0ZXN0dG9rZW4="},
				AdapterGetPageFuncs: make(map[string]internal.AdapterGetPageFunc),
			},
		},
		"no_tokens_at_path": {
			inputAuthTokensPath: "/",
			inputStopChan:       nil,
			wantAdapterServer: &internal.Server{
				Tokens:              nil,
				AdapterGetPageFuncs: make(map[string]internal.AdapterGetPageFunc),
			},
		},
		"invalid_tokens_at_path": {
			inputAuthTokensPath: invalidTokensPath,
			inputStopChan:       nil,
			wantAdapterServer: &internal.Server{
				Tokens:              nil,
				AdapterGetPageFuncs: make(map[string]internal.AdapterGetPageFunc),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotAdapterServer := newWithAuthTokensPath(
				tc.inputAuthTokensPath,
				tc.inputStopChan,
			)

			AssertDeepEqual(t, tc.wantAdapterServer, gotAdapterServer)
		})
	}
}

func TestNewWithAuthTokensPathFileWatcher(t *testing.T) {
	validTokensPath := "./TOKENS_1"

	tokens := []byte(`["dGhpc2lzYXRlc3R0b2tlbg==","dGhpc2lzYWxzb2F0ZXN0dG9rZW4="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})

	gotAdapterServer := newWithAuthTokensPath(
		validTokensPath,
		stop,
	)

	// Assert the initial state of the tokens are correct
	AssertDeepEqual(t, gotAdapterServer.(*internal.Server).Tokens, []string{"dGhpc2lzYXRlc3R0b2tlbg==", "dGhpc2lzYWxzb2F0ZXN0dG9rZW4="})

	// Add a third token to the file
	tokens = []byte(`["dGhpc2lzYXRlc3R0b2tlbg==","dGhpc2lzYWxzb2F0ZXN0dG9rZW4=","TfGX4vJkrqfRyvUviDpj3Q=="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Assert the tokens have been updated
	AssertDeepEqual(t, gotAdapterServer.(*internal.Server).Tokens, []string{
		"dGhpc2lzYXRlc3R0b2tlbg==", "dGhpc2lzYWxzb2F0ZXN0dG9rZW4=", "TfGX4vJkrqfRyvUviDpj3Q==",
	})

	// Close the file watcher using the stop channel
	close(stop)

	// Remove the first 2 tokens from the file
	tokens = []byte(`"TfGX4vJkrqfRyvUviDpj3Q=="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)

	// Assert the tokens have not been updated (e.g. that the file watcher was closed correctly)
	AssertDeepEqual(t, gotAdapterServer.(*internal.Server).Tokens, []string{
		"dGhpc2lzYXRlc3R0b2tlbg==", "dGhpc2lzYWxzb2F0ZXN0dG9rZW4=", "TfGX4vJkrqfRyvUviDpj3Q==",
	})
}
