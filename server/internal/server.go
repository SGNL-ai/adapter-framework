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
	"context"
	"fmt"
	"sync"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/grpc/codes"
	grpc_metadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Server is an implementation of the AdapterServer gRPC service which
// delegates the implementation of the RPCs to high-level Adapter
// implementation based on a provided type, and translates and
// validates RPC requests and responses.
//
// The Config type parameter must be a struct type into which the configuration
// JSON object can be unmarshaled into.
type Server[Config any] struct {
	api_adapter_v1.UnimplementedAdapterServer

	// Adapters contains a map of high-level implementations of the service
	// as well as their associated types.
	// The key in this map should match the Supported Datasource Type
	// specified on the Adapter object created in SGNL.
	Adapters map[string]framework.Adapter[Config]

	// Tokens contains a lists of valid auth tokens for this server. This list of Tokens
	// is populated when the server is created based on the JSON-encoded value in the file
	// located under the path contained in the `AUTH_TOKENS_PATH` environment variable and is
	// updated any time this file is modified.
	// This field must only be accessed for reading or writing while locking TokensMutex.
	Tokens []string

	// TokensMutex is the mutex that must be locked for every access to Tokens.
	TokensMutex sync.RWMutex
}

func (s *Server[Config]) GetPage(ctx context.Context, req *api_adapter_v1.GetPageRequest) (*api_adapter_v1.GetPageResponse, error) {
	if err := s.validateAuthenticationToken(ctx); err != nil {
		return nil, err
	}

	adapterRequest, reverseMapping, adapterErr := getAdapterRequest[Config](req)

	if adapterErr != nil {
		return api_adapter_v1.NewGetPageResponseError(adapterErr), nil
	}

	if adapter, ok := s.Adapters[req.Datasource.Type]; ok {
		adapterResponse := adapter.GetPage(ctx, adapterRequest)

		return getResponse(reverseMapping, &adapterResponse), nil
	}

	adapterErr = &api_adapter_v1.Error{
		Message: fmt.Sprintf("Unsupported datasource type provided: %s.", req.Datasource.Type),
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
	}

	return api_adapter_v1.NewGetPageResponseError(adapterErr), nil
}

// validateAuthenticationToken verifies the request has the correct token to access the
// adapter. Will return nil if the provided token matches any of the tokens
// specified in the file located at AUTH_TOKENS_PATH.
// Otherwise, will return an error.
func (s *Server[Config]) validateAuthenticationToken(ctx context.Context) error {
	metadata, ok := grpc_metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "invalid or missing token")
	}

	requestTokens := metadata.Get("token")
	if len(requestTokens) != 1 {
		return status.Errorf(codes.Unauthenticated, "invalid or missing token")
	}

	s.TokensMutex.RLock()
	defer s.TokensMutex.RUnlock()

	// TODO: After upgrading go to 1.21+, replace with the `Contains` method
	for _, y := range s.Tokens {
		if y == requestTokens[0] {
			return nil
		}
	}

	return status.Errorf(codes.Unauthenticated, "invalid or missing token")
}
