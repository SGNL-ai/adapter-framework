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

import "sync"

// MockLogger is a test implementation of the Logger interface.
// It records all log calls for verification in tests.
type MockLogger struct {
	mu      sync.Mutex
	entries []LogEntry
	fields  []Field
}

var _ Logger = (*MockLogger)(nil)

// LogEntry represents a single log entry.
type LogEntry struct {
	Level   string
	Message string
	Fields  []Field
}

// NewMockLogger creates a new MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		entries: make([]LogEntry, 0),
		fields:  make([]Field, 0),
	}
}

// Info logs an informational message.
func (m *MockLogger) Info(msg string, fields ...Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, LogEntry{
		Level:   "info",
		Message: msg,
		Fields:  append(m.fields, fields...),
	})
}

// Error logs an error message.
func (m *MockLogger) Error(msg string, fields ...Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, LogEntry{
		Level:   "error",
		Message: msg,
		Fields:  append(m.fields, fields...),
	})
}

// Debug logs a debug message.
func (m *MockLogger) Debug(msg string, fields ...Field) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, LogEntry{
		Level:   "debug",
		Message: msg,
		Fields:  append(m.fields, fields...),
	})
}

// With creates a child logger with pre-attached fields.
func (m *MockLogger) With(fields ...Field) Logger {
	m.mu.Lock()
	defer m.mu.Unlock()

	return &MockLogger{
		entries: m.entries,
		fields:  append(append([]Field{}, m.fields...), fields...),
	}
}

// Entries returns all logged entries.
func (m *MockLogger) Entries() []LogEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]LogEntry{}, m.entries...)
}

// Fields returns all pre-attached fields.
func (m *MockLogger) Fields() []Field {
	m.mu.Lock()
	defer m.mu.Unlock()

	return append([]Field{}, m.fields...)
}

// Reset clears all logged entries.
func (m *MockLogger) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = make([]LogEntry, 0)
}
