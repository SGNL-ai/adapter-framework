package server

import (
	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	"github.com/sgnl-ai/adapter-framework/server/internal"
)

// New returns an AdapterServer that wraps the given high-level
// Adapter implementation.
func New[Config any](adapter framework.Adapter[Config]) api_adapter_v1.AdapterServer {
	return &internal.Server[Config]{
		Adapter: adapter,
	}
}
