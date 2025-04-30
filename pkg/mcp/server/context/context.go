package context

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type CtxRequestIdKey struct{}
type CtxSessionIdKey struct{}
type CtxRootSpanKey struct{}
type CtxStartTimeKey struct{}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, CtxRequestIdKey{}, requestID)
}

func CtxRequestId(ctx context.Context) string {
	if rid, ok := ctx.Value(CtxRequestIdKey{}).(string); ok {
		return rid
	}

	return ""
}

func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, CtxSessionIdKey{}, sessionID)
}

func CtxSessionID(ctx context.Context) string {
	if sid, ok := ctx.Value(CtxSessionIdKey{}).(string); ok {
		return sid
	}

	return ""
}

func WithRootSpan(ctx context.Context, span trace.Span) context.Context {
	return context.WithValue(ctx, CtxRootSpanKey{}, span)
}

func CtxRootSpan(ctx context.Context) trace.Span {
	if span, ok := ctx.Value(CtxRootSpanKey{}).(trace.Span); ok {
		return span
	}

	return trace.SpanFromContext(ctx)
}

func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, CtxStartTimeKey{}, t)
}

func CtxStartTime(ctx context.Context) time.Time {
	if t, ok := ctx.Value(CtxStartTimeKey{}).(time.Time); ok {
		return t
	}

	return time.Now()
}
