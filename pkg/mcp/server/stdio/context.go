package stdio

import (
	"context"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	System    = "mcpserver"
	Transport = "stdio"
)

type CtxRootSpanKey struct{}

type MCPStdioServerContextHandler struct {
	generator      uuid.UuidGenerator
	tracerProvider oteltrace.TracerProvider
	logger         *log.Logger
}

func NewMCPStdioServerContextHandler(
	generator uuid.UuidGenerator,
	tracerProvider oteltrace.TracerProvider,
	logger *log.Logger,
) *MCPStdioServerContextHandler {
	return &MCPStdioServerContextHandler{
		generator:      generator,
		tracerProvider: tracerProvider,
		logger:         logger,
	}
}

func (h *MCPStdioServerContextHandler) Handle() server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
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
			),
		)

		ctx = context.WithValue(ctx, CtxRootSpanKey{}, span)

		// logger propagation
		logger := h.logger.
			With().
			Str("system", System).
			Str("transport", Transport).
			Logger()

		logger.
			Info().
			Msg("MCP Stdio request")

		return logger.WithContext(ctx)
	}
}
