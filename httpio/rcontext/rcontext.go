// This package provides support for managing request and response context information.
package rcontext

import (
	"context"
	"time"
)

type contextKey string

const requestInfoKey = contextKey("requestInfo")
const ResponseInfoKey = contextKey("ResponseInfo")

// RequestInfo holds information about the incoming request.
type RequestInfo struct {
	RequestID string
	StartedAt time.Time
}

// ResponseInfo holds information about the outgoing response.
type ResponseInfo struct {
	StatusCode int
}

// GetRequestInfo returns the RequestInfo from the context.
func GetRequestInfo(ctx context.Context) *RequestInfo {
	info, ok := ctx.Value(requestInfoKey).(*RequestInfo)
	if !ok {
		return nil
	}

	return info
}

// SetRequestInfo sets the RequestInfo in the context.
func SetRequestInfo(ctx context.Context, info *RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, info)
}

// GetResponseInfo returns the ResponseInfo from the context.
func GetResponseInfo(ctx context.Context) *ResponseInfo {
	info, ok := ctx.Value(ResponseInfoKey).(*ResponseInfo)
	if !ok {
		return nil
	}

	return info
}

// SetResponseInfo sets the ResponseInfo in the context.
func SetResponseInfo(ctx context.Context, info *ResponseInfo) context.Context {
	return context.WithValue(ctx, ResponseInfoKey, info)
}
