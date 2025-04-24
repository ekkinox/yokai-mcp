package sse

import (
	"context"
	"net/http"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	System          = "mcpserver"
	Transport       = "sse"
	HeaderRequestID = "X-Request-Id"
	RequestID       = "requestID"
	SessionID       = "sessionID"
)

type CtxRequestIdKey struct{}
type CtxSessionIdKey struct{}
type CtxRootSpanKey struct{}

type MCPSSEServerContextHandler struct {
	generator      uuid.UuidGenerator
	tracerProvider oteltrace.TracerProvider
	logger         *log.Logger
}

func NewMCPSSEServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider oteltrace.TracerProvider,
	logger *log.Logger,
) *MCPSSEServerContextHandler {
	return &MCPSSEServerContextHandler{
		generator:      generator,
		tracerProvider: tracerProvider,
		logger:         logger,
	}
}

func (h *MCPSSEServerContextHandler) Handle() server.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		// sessionId propagation
		sID := r.URL.Query().Get("sessionId")

		ctx = context.WithValue(ctx, CtxSessionIdKey{}, sID)

		// requestId propagation
		rID := r.Header.Get(HeaderRequestID)

		if rID == "" {
			rID = h.generator.Generate()
			r.Header.Set(HeaderRequestID, rID)
		}

		ctx = context.WithValue(ctx, CtxRequestIdKey{}, rID)

		// tracer propagation
		ctx = trace.WithContext(ctx, h.tracerProvider)

		ctx, span := trace.CtxTracer(ctx).Start(
			ctx,
			"MCP",
			oteltrace.WithNewRoot(),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(
				attribute.String("system", System),
				attribute.String("transport", Transport),
				attribute.String(SessionID, sID),
				attribute.String(RequestID, rID),
			),
		)

		ctx = context.WithValue(ctx, CtxRootSpanKey{}, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", System).
			Str("transport", Transport).
			Str(SessionID, sID).
			Str(RequestID, rID).
			Logger()

		logger.
			Info().
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("userAgent", r.UserAgent()).
			Msg("MCP SSE request")

		return logger.WithContext(ctx)
	}
}
