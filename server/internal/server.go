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

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
)

// Server is an implementation of the AdapterServer gRPC service which
// delegates the implementation of the RPCs to a high-level Adapter
// implementation, and translates and validates RPC requests and responses.
//
// The Config type parameter must be a struct type into which the configuration
// JSON object can be unmarshaled into.
type Server[Config any] struct {
	api_adapter_v1.UnimplementedAdapterServer

	// Adapter is the high-level implementation of the service.
	Adapter framework.Adapter[Config]
}

func (s *Server[Config]) GetPage(ctx context.Context, req *api_adapter_v1.GetPageRequest) (*api_adapter_v1.GetPageResponse, error) {
	adapterRequest, reverseMapping, adapterErr := getAdapterRequest[Config](req)

	if adapterErr != nil {
		return api_adapter_v1.NewGetPageResponseError(adapterErr), nil
	}

	adapterResponse := s.Adapter.GetPage(ctx, adapterRequest)

	return getResponse(reverseMapping, &adapterResponse), nil
}
