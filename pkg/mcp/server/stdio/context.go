package stdio

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	yokaimcpservercontext "github.com/ekkinox/yokai-mcp/pkg/mcp/server/context"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ MCPStdioServerContextHandler = (*DefaultMCPStdioServerContextHandler)(nil)

type MCPStdioServerContextHandler interface {
	Handle() server.StdioContextFunc
}

type DefaultMCPStdioServerContextHandler struct {
	generator      uuid.UuidGenerator
	tracerProvider oteltrace.TracerProvider
	logger         *log.Logger
}

func NewDefaultMCPStdioServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider oteltrace.TracerProvider,
	logger *log.Logger,
) *DefaultMCPStdioServerContextHandler {
	return &DefaultMCPStdioServerContextHandler{
		generator:      generator,
		tracerProvider: tracerProvider,
		logger:         logger,
	}
}

func (h *DefaultMCPStdioServerContextHandler) Handle() server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
		// start time propagation
		ctx = yokaimcpservercontext.WithStartTime(ctx, time.Now())

		// requestId propagation
		rID := h.generator.Generate()

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
				attribute.String("mcp.transport", "stdio"),
				attribute.String("mcp.requestID", rID),
			),
		)

		ctx = yokaimcpservercontext.WithRootSpan(ctx, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", "mcpserver").
			Str("mcpTransport", "stdio").
			Str("mcpRequestID", rID).
			Logger()

		return logger.WithContext(ctx)
	}
}
