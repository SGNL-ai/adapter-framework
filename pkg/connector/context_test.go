package connector

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantVal ConnectorInfo
		wantOk  bool
	}{
		{
			name:    "empty context returns false",
			ctx:     context.Background(),
			wantVal: ConnectorInfo{},
			wantOk:  false,
		},
		{
			name: "context with value returns true",
			ctx: context.WithValue(context.Background(), key{}, ConnectorInfo{
				ID: "test-id",
			}),
			wantVal: ConnectorInfo{ID: "test-id"},
			wantOk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := FromContext(tt.ctx)
			if gotOk != tt.wantOk {
				t.Errorf("FromContext() ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotVal != tt.wantVal {
				t.Errorf("FromContext() val = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}

func TestWithContext(t *testing.T) {
	tests := []struct {
		name        string
		parent      context.Context
		val         ConnectorInfo
		wantErr     bool
		wantErrText string
	}{
		{
			name:    "successfully add context info",
			parent:  context.Background(),
			val:     ConnectorInfo{ID: "test-id"},
			wantErr: false,
		},
		{
			name:        "error when context already has value",
			parent:      context.WithValue(context.Background(), key{}, ConnectorInfo{ID: "existing-id"}),
			val:         ConnectorInfo{ID: "new-id"},
			wantErr:     true,
			wantErrText: "context is already configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := WithContext(tt.parent, tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("WithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err == nil || !contains(err.Error(), tt.wantErrText) {
					t.Errorf("WithContext() error = %v, want error containing %v", err, tt.wantErrText)
				}
				return
			}

			// Verify the context value was set correctly
			gotVal, ok := FromContext(ctx)
			if !ok {
				t.Error("WithContext() value not found in context")
			}
			if gotVal != tt.val {
				t.Errorf("WithContext() value = %v, want %v", gotVal, tt.val)
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
