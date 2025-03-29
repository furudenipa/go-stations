package context

import "context"

type osContextKey string

func WithOS(ctx context.Context, os string) context.Context {
	return context.WithValue(ctx, osContextKey("os"), os)
}

func OS(ctx context.Context) string {
	if os, ok := ctx.Value(osContextKey("os")).(string); ok {
		return os
	}
	return ""
}
