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
	"context"
	"errors"
	"os"
	"sort"
	"testing"
	"time"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/pkg/logs/zaplog"
	"github.com/sgnl-ai/adapter-framework/server/internal"
	"go.uber.org/zap"
)

type MockAdapterA struct{}

func (a *MockAdapterA) GetPage(ctx context.Context, request *framework.Request[TestConfigA]) framework.Response {
	return framework.Response{}
}

func NewAdapterA() framework.Adapter[TestConfigA] {
	return &MockAdapterA{}
}

type MockAdapterB struct{}

func (a *MockAdapterB) GetPage(ctx context.Context, request *framework.Request[TestConfigB]) framework.Response {
	return framework.Response{}
}

func NewAdapterB() framework.Adapter[TestConfigB] {
	return &MockAdapterB{}
}

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
				nil,
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
		nil,
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

func TestRegisterAdapter(t *testing.T) {
	s := &internal.Server{
		AdapterGetPageFuncs: make(map[string]internal.AdapterGetPageFunc),
	}

	var adapterServer api_adapter_v1.AdapterServer = s

	if err := RegisterAdapter(adapterServer, "Mock-1.0.1", NewAdapterA()); err != nil {
		t.Fatal(err)
	}

	if err := RegisterAdapter(adapterServer, "Mock-1.0.2", NewAdapterB()); err != nil {
		t.Fatal(err)
	}

	var registeredDatasources []string

	for k := range s.AdapterGetPageFuncs {
		registeredDatasources = append(registeredDatasources, k)
	}

	sort.Strings(registeredDatasources)

	want := []string{"Mock-1.0.1", "Mock-1.0.2"}

	AssertDeepEqual(t, want, registeredDatasources)
}

func TestRegisterAdapterInvalidServer(t *testing.T) {
	type InvalidServer struct {
		api_adapter_v1.UnimplementedAdapterServer
	}

	s := &InvalidServer{}

	var adapterServer api_adapter_v1.AdapterServer = s

	err := RegisterAdapter(adapterServer, "Mock-1.0.1", NewAdapterA())

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	AssertDeepEqual(t, err, errors.New("type assertion to *internal.Server failed"))
}

func TestNew_WithLogger(t *testing.T) {
	validTokensPath := "./TOKENS_WITH_LOGGER"

	tokens := []byte(`["dGhpc2lzYXRlc3R0b2tlbg=="]`)
	if err := os.WriteFile(validTokensPath, tokens, 0666); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validTokensPath)

	t.Setenv("AUTH_TOKENS_PATH", validTokensPath)

	zapLogger := zap.NewNop()
	logger := zaplog.New(zapLogger)
	stop := make(chan struct{})
	defer close(stop)

	server := New(stop, WithLogger(logger))

	internalServer, ok := server.(*internal.Server)
	if !ok {
		t.Fatal("Expected *internal.Server")
	}

	if internalServer.Logger == nil {
		t.Error("Expected logger to be set")
	}
}
