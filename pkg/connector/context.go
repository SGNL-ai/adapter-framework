// Copyright 2025 SGNL.ai, Inc.

package connector

import (
	"context"
	"fmt"
)

// SourceType identifies a particular source as a datasource or
// an integration.
type SourceType int

const (
	Unknown SourceType = iota
	Datasource
	Integration
)

// ConnectorInfo carries information about the on-premises connector,
// such as its ID, and the source it belongs to, for communication purposes.
type ConnectorInfo struct {
	ID         string
	ClientID   string
	TenantID   string
	SourceType SourceType
	SourceID   string
}

// key for storing the ConnectorInfo in a derived context.
type key struct{}

// FromContext returns the ConnectorInfo value stored in the ctx, if any.
func FromContext(ctx context.Context) (v ConnectorInfo, ok bool) {
	if ctx == nil {
		panic("cannot read connector info value from nil context")
	}
	v, ok = ctx.Value(key{}).(ConnectorInfo)
	return
}

// WithContext returns a derived ctx with the ConnectorInfo value stored in it.
func WithContext(parent context.Context, val ConnectorInfo) (context.Context, error) {
	if v, ok := FromContext(parent); ok {
		return nil, fmt.Errorf("context is already configured with the connector info, %v", v)
	}
	return context.WithValue(parent, key{}, val), nil
}
