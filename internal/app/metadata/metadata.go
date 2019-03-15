package metadata

import (
	"context"
)

type metaKey struct{}

// Metadata is our way of representing incoming request headers internally.
type Metadata map[string]string

func FromContext(ctx context.Context) (Metadata, bool) {
	md, ok := ctx.Value(metaKey{}).(Metadata)
	return md, ok
}

func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metaKey{}, md)
}

func GetRequestID(ctx context.Context) string {
	if md, ok := FromContext(ctx); ok {
		if val, ok := md["x-request-id"]; ok {
			return val
		}
	}
	return ""
}
