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

// Logger is a minimal interface for structured logging.
type Logger interface {
	// Info logs an informational message with optional fields.
	Info(msg string, fields ...Field)

	// Error logs an error message with optional fields.
	Error(msg string, fields ...Field)

	// Debug logs a debug message with optional fields.
	Debug(msg string, fields ...Field)

	// With creates a child logger with the given fields pre-attached.
	With(fields ...Field) Logger
}

// Field represents a key-value pair for structured logging.
type Field struct {
	Key   string
	Value any
}
