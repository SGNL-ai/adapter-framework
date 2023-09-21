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
	"encoding/json"
	"fmt"
	"os"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"google.golang.org/grpc/codes"
	grpc_metadata "google.golang.org/grpc/metadata"
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
}

func (s *Server[Config]) GetPage(ctx context.Context, req *api_adapter_v1.GetPageRequest) (*api_adapter_v1.GetPageResponse, error) {
	if err := validateAuthenticationToken(ctx); err != nil {
		return api_adapter_v1.NewGetPageResponseError(err), nil
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
func validateAuthenticationToken(ctx context.Context) *api_adapter_v1.Error {
	metadata, ok := grpc_metadata.FromIncomingContext(ctx)
	if !ok {
		return &api_adapter_v1.Error{
			Message: "Invalid or missing token.",
			Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
		}
	}

	requestTokens := metadata.Get("token")
	if len(requestTokens) != 1 {
		return &api_adapter_v1.Error{
			Message: "Invalid or missing token.",
			Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
		}
	}

	path, exists := os.LookupEnv("AUTH_TOKENS_PATH")
	if !exists {
		return &api_adapter_v1.Error{
			Message: "Invalid or missing token.",
			Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
		}
	}

	jsonValidTokens, err := os.ReadFile(path)
	if err != nil {
		return &api_adapter_v1.Error{
			Message: "Invalid or missing token.",
			Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
		}
	}

	validTokens := new([]string)

	if err := json.Unmarshal(jsonValidTokens, validTokens); err != nil || validTokens == nil {
		return &api_adapter_v1.Error{
			Message: "Invalid or missing token.",
			Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
		}
	}

	// TODO: After upgrading go to 1.21+, replace with the `Contains` method
	for _, y := range *validTokens {
		if y == requestTokens[0] {
			return nil
		}
	}

	return &api_adapter_v1.Error{
		Message: "Invalid or missing token.",
		Code:    api_adapter_v1.ErrorCode(codes.Unauthenticated),
	}
}
