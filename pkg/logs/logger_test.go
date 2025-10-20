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

package logs

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestContextWithLogger(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	newCtx := ContextWithLogger(ctx, logger)

	if newCtx == nil {
		t.Fatal("ContextWithLogger returned nil context")
	}

	if newCtx == ctx {
		t.Error("ContextWithLogger should return a new context, not the same one")
	}
}

func TestLoggerFromContext(t *testing.T) {
	logger := zap.NewNop()

	tests := map[string]struct {
		setupCtx   func() context.Context
		wantLogger *zap.Logger
	}{
		"returns_logger_when_present": {
			setupCtx: func() context.Context {
				ctx := context.Background()

				return ContextWithLogger(ctx, logger)
			},
			wantLogger: logger,
		},
		"returns_nil_when_logger_not_present": {
			setupCtx: func() context.Context {
				return context.Background()
			},
			wantLogger: nil,
		},
		"returns_nil_when_context_has_wrong_type": {
			setupCtx: func() context.Context {
				ctx := context.Background()

				return context.WithValue(ctx, loggerContextKey{}, "not a logger")
			},
			wantLogger: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := tc.setupCtx()
			retrievedLogger := LoggerFromContext(ctx)

			if tc.wantLogger != retrievedLogger {
				t.Error("LoggerFromContext returned different logger than the one stored")
			}
		})
	}
}
