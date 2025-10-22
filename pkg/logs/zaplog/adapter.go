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

package zaplog

import (
	"github.com/sgnl-ai/adapter-framework/pkg/logs"
	"go.uber.org/zap"
)

// Adapter wraps a *zap.Logger to implement the logs.Logger interface.
type Adapter struct {
	logger *zap.Logger
}

// New creates a new logs.Logger from a *zap.Logger.
func New(logger *zap.Logger) logs.Logger {
	return &Adapter{logger: logger}
}

// Info logs an informational message.
func (a *Adapter) Info(msg string, fields ...logs.Field) {
	a.logger.Info(msg, toZapFields(fields)...)
}

// Error logs an error message.
func (a *Adapter) Error(msg string, fields ...logs.Field) {
	a.logger.Error(msg, toZapFields(fields)...)
}

// Debug logs a debug message.
func (a *Adapter) Debug(msg string, fields ...logs.Field) {
	a.logger.Debug(msg, toZapFields(fields)...)
}

// With creates a child logger with pre-attached fields.
func (a *Adapter) With(fields ...logs.Field) logs.Logger {
	return &Adapter{
		logger: a.logger.With(toZapFields(fields)...),
	}
}

// Unwrap returns the underlying *zap.Logger.
// This allows consumers to access zap-specific features when needed.
func (a *Adapter) Unwrap() *zap.Logger {
	return a.logger
}

// UnwrapLogger attempts to extract a *zap.Logger from a logs.Logger.
// Returns the underlying *zap.Logger and true if the logger is a zaplog.Adapter,
// otherwise returns nil and false.
//
// Example usage:
//
//	logger := logs.FromContext(ctx)
//	if zapLogger, ok := zaplog.UnwrapLogger(logger); ok {
//	    // Use zap-specific features
//	    zapLogger.Sync()
//	}
func UnwrapLogger(logger logs.Logger) (*zap.Logger, bool) {
	if adapter, ok := logger.(*Adapter); ok {
		return adapter.Unwrap(), true
	}
	return nil, false
}

// toZapFields converts logs.Field to zap.Field.
func toZapFields(fields []logs.Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
