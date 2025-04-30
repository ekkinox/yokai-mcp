package sse

import (
	"context"
	"net/http"
	"time"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	yokaimcpservercontext "github.com/ekkinox/yokai-mcp/pkg/mcp/server/context"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ MCPSSEServerContextHandler = (*DefaultMCPSSEServerContextHandler)(nil)

type MCPSSEServerContextHandler interface {
	Handle() server.SSEContextFunc
}

type DefaultMCPSSEServerContextHandler struct {
	generator      uuid.UuidGenerator
	tracerProvider oteltrace.TracerProvider
	logger         *log.Logger
}

func NewDefaultMCPSSEServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider oteltrace.TracerProvider,
	logger *log.Logger,
) *DefaultMCPSSEServerContextHandler {
	return &DefaultMCPSSEServerContextHandler{
		generator:      generator,
		tracerProvider: tracerProvider,
		logger:         logger,
	}
}

func (h *DefaultMCPSSEServerContextHandler) Handle() server.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		// start time propagation
		ctx = yokaimcpservercontext.WithStartTime(ctx, time.Now())

		// sessionId propagation
		sID := r.URL.Query().Get("sessionId")

		ctx = yokaimcpservercontext.WithSessionID(ctx, sID)

		// requestId propagation
		rID := r.Header.Get("X-Request-Id")

		if rID == "" {
			rID = h.generator.Generate()
			r.Header.Set("X-Request-Id", rID)
		}

		ctx = yokaimcpservercontext.WithRequestID(ctx, rID)

		// tracer propagation
		ctx = trace.WithContext(ctx, h.tracerProvider)

		ctx, span := trace.CtxTracer(ctx).Start(
			ctx,
			"MCP",
			oteltrace.WithNewRoot(),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(
				attribute.String("system", "mcpserver"),
				attribute.String("mcp.transport", "sse"),
				attribute.String("mcp.sessionID", sID),
				attribute.String("mcp.requestID", rID),
			),
		)

		ctx = yokaimcpservercontext.WithRootSpan(ctx, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", "mcpserver").
			Str("mcpTransport", "sse").
			Str("mcpSessionID", sID).
			Str("mcpRequestID", rID).
			Logger()

		return logger.WithContext(ctx)
	}
}
