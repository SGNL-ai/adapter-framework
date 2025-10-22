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
	"log/slog"

	"github.com/sgnl-ai/adapter-framework/pkg/logs"
)

// Adapter wraps a *slog.Logger to implement the logs.Logger interface.
type Adapter struct {
	logger *slog.Logger
}

// New creates a new logs.Logger from a *slog.Logger.
func New(logger *slog.Logger) logs.Logger {
	return &Adapter{logger: logger}
}

// Info logs an informational message.
func (a *Adapter) Info(msg string, fields ...logs.Field) {
	a.logger.Info(msg, toSlogArgs(fields)...)
}

// Error logs an error message.
func (a *Adapter) Error(msg string, fields ...logs.Field) {
	a.logger.Error(msg, toSlogArgs(fields)...)
}

// Debug logs a debug message.
func (a *Adapter) Debug(msg string, fields ...logs.Field) {
	a.logger.Debug(msg, toSlogArgs(fields)...)
}

// With creates a child logger with pre-attached fields.
func (a *Adapter) With(fields ...logs.Field) logs.Logger {
	return &Adapter{
		logger: a.logger.With(toSlogArgs(fields)...),
	}
}

// Unwrap returns the underlying *slog.Logger.
// This allows consumers to access slog-specific features when needed.
func (a *Adapter) Unwrap() *slog.Logger {
	return a.logger
}

// UnwrapLogger attempts to extract a *slog.Logger from a logs.Logger.
// Returns the underlying *slog.Logger and true if the logger is a sloglog.Adapter,
// otherwise returns nil and false.
//
// Example usage:
//
//	logger := logs.FromContext(ctx)
//	if slogLogger, ok := sloglog.UnwrapLogger(logger); ok {
//	    // Use slog-specific features
//	    slogLogger.With("key", "value").Info("message")
//	}
func UnwrapLogger(logger logs.Logger) (*slog.Logger, bool) {
	if adapter, ok := logger.(*Adapter); ok {
		return adapter.Unwrap(), true
	}
	return nil, false
}

// toSlogArgs converts logs.Field to slog arguments.
func toSlogArgs(fields []logs.Field) []any {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return args
}
