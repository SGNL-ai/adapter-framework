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

package sloglog

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/sgnl-ai/adapter-framework/pkg/logs"
)

func TestSlogLogger_Implementation(t *testing.T) {
	// Create an slog logger
	slogLogger := slog.Default()

	// Wrap it in our interface
	logger := New(slogLogger)

	// Test that it implements the Logger interface
	var _ logs.Logger = logger

	// Test basic logging operations
	logger.Info("test message")
	logger.Error("error message")
	logger.Debug("debug message")

	// Test With method
	childLogger := logger.With(logs.ClientID("test-client"), logs.TenantID("test-tenant"))
	childLogger.Info("child logger message")
}

func TestSlogLogger_Fields(t *testing.T) {
	// Create a test logger with a buffer
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slogLogger := slog.New(handler)

	logger := New(slogLogger)

	// Log with fields
	logger.With(
		logs.ClientID("client-123"),
		logs.TenantID("tenant-456"),
	).Info("test message")

	// Verify output contains the fields (basic check)
	output := buf.String()
	if output == "" {
		t.Error("Expected log output, got empty string")
	}
}

func TestUnwrapLogger(t *testing.T) {
	slogLogger := slog.Default()
	logger := New(slogLogger)

	// Test successful unwrap
	unwrapped, ok := UnwrapLogger(logger)
	if !ok {
		t.Fatal("Expected UnwrapLogger to return true for sloglog.Adapter")
	}
	if unwrapped != slogLogger {
		t.Error("Expected unwrapped logger to be the same as the original slog logger")
	}

	// Test unwrap with child logger
	childLogger := logger.With(logs.ClientID("test"))
	unwrappedChild, ok := UnwrapLogger(childLogger)
	if !ok {
		t.Fatal("Expected UnwrapLogger to work with child logger")
	}
	if unwrappedChild == nil {
		t.Error("Expected unwrapped child logger to be non-nil")
	}
}

func TestAdapter_Unwrap(t *testing.T) {
	slogLogger := slog.Default()
	adapter := New(slogLogger).(*Adapter)

	unwrapped := adapter.Unwrap()
	if unwrapped != slogLogger {
		t.Error("Expected Unwrap to return the original slog logger")
	}
}
