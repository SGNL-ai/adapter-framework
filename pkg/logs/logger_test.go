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

package logs_test

import (
	"context"
	"testing"

	"github.com/sgnl-ai/adapter-framework/pkg/logs"
	"github.com/sgnl-ai/adapter-framework/pkg/logs/zaplog"
	"go.uber.org/zap"
)

func TestContextWithLogger(t *testing.T) {
	zapLogger := zap.NewNop()
	logger := zaplog.New(zapLogger)
	ctx := context.Background()

	newCtx := logs.NewContextWithLogger(ctx, logger)

	if newCtx == nil {
		t.Fatal("ContextWithLogger returned nil context")
	}

	if newCtx == ctx {
		t.Error("ContextWithLogger should return a new context, not the same one")
	}
}

func TestLoggerFromContext(t *testing.T) {
	zapLogger := zap.NewNop()
	logger := zaplog.New(zapLogger)

	tests := map[string]struct {
		setupCtx   func() context.Context
		wantLogger bool
	}{
		"returns_logger_when_present": {
			setupCtx: func() context.Context {
				ctx := context.Background()
				return logs.NewContextWithLogger(ctx, logger)
			},
			wantLogger: true,
		},
		"returns_nil_when_logger_not_present": {
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantLogger: false,
		},
		"returns_nil_when_context_has_wrong_type": {
			setupCtx: func() context.Context {
				ctx := context.Background()
				type loggerContextKey struct{}
				return context.WithValue(ctx, loggerContextKey{}, "not a logger")
			},
			wantLogger: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := tc.setupCtx()
			retrievedLogger := logs.FromContext(ctx)

			if tc.wantLogger {
				if retrievedLogger == nil {
					t.Fatal("LoggerFromContext returned nil, expected logger")
				}
			} else {
				if retrievedLogger != nil {
					t.Error("LoggerFromContext should return nil")
				}
			}
		})
	}
}
