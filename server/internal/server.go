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

type AdapterGetPageFunc func(ctx context.Context, req *api_adapter_v1.GetPageRequest) (framework.Response, *EntityReverseIdMapping)

// Server is an implementation of the AdapterServer gRPC service which
// delegates the implementation of the RPCs to high-level Adapter
// implementation based on a provided type, and translates and
// validates RPC requests and responses.
//
// The Config type parameter must be a struct type into which the configuration
// JSON object can be unmarshaled into.
type Server struct {
	api_adapter_v1.UnimplementedAdapterServer

	// TODO
	AdapterGetPageFuncs map[string]AdapterGetPageFunc

	// Tokens contains a lists of valid auth tokens for this server. This list of Tokens
	// is populated when the server is created based on the JSON-encoded value in the file
	// located under the path contained in the `AUTH_TOKENS_PATH` environment variable and is
	// updated any time this file is modified.
	// This field must only be accessed for reading or writing while locking TokensMutex.
	Tokens []string

	// TokensMutex is the mutex that must be locked for every access to Tokens.
	TokensMutex sync.RWMutex
}

func (s *Server) GetPage(ctx context.Context, req *api_adapter_v1.GetPageRequest) (*api_adapter_v1.GetPageResponse, error) {
	if err := s.validateAuthenticationToken(ctx); err != nil {
		return nil, err
	}

	if adapterGetPageFunc, ok := s.AdapterGetPageFuncs[req.Datasource.Type]; ok {
		adapterResponse, reverseMapping := adapterGetPageFunc(ctx, req)

		return getResponse(reverseMapping, &adapterResponse), nil
	}

	adapterErr := &api_adapter_v1.Error{
		Message: fmt.Sprintf("Unsupported datasource type provided: %s.", req.Datasource.Type),
		Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INVALID_DATASOURCE_CONFIG,
	}

	return api_adapter_v1.NewGetPageResponseError(adapterErr), nil
}

// validateAuthenticationToken verifies the request has the correct token to access the
// adapter. Will return nil if the provided token matches any of the tokens
// specified in the file located at AUTH_TOKENS_PATH.
// Otherwise, will return an error.
func (s *Server) validateAuthenticationToken(ctx context.Context) error {
	metadata, ok := grpc_metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "invalid or missing token")
	}

	requestTokens := metadata.Get("token")
	if len(requestTokens) != 1 {
		return status.Error(codes.Unauthenticated, "invalid or missing token")
	}

	s.TokensMutex.RLock()
	defer s.TokensMutex.RUnlock()

	// TODO: After upgrading go to 1.21+, replace with the `Contains` method
	for _, y := range s.Tokens {
		if y == requestTokens[0] {
			return nil
		}
	}

	return status.Error(codes.Unauthenticated, "invalid or missing token")
}

func RegisterAdapter[Config any](s *Server, datasourceType string, adapter framework.Adapter[Config]) error {
	// TODO: check for duplicates and return an error

	s.AdapterGetPageFuncs[datasourceType] = func(ctx context.Context, req *api_adapter_v1.GetPageRequest) (framework.Response, *EntityReverseIdMapping) {
		adapterRequest, reverseMapping, adapterErr := GetAdapterRequest[Config](req)
		if adapterErr != nil {
			return framework.NewGetPageResponseError(&framework.Error{
				Message: adapterErr.Message,
				Code:    adapterErr.Code,
				// TODO
				// RetryAfter: adapterErr.RetryAfter,
			}), nil
		}

		return adapter.GetPage(ctx, adapterRequest), reverseMapping
	}

	return nil
}
