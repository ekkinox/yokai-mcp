package mcp

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
	HeaderRequestID = "X-Request-Id"
	RequestID       = "requestID"
	SessionID       = "sessionID"
	Transport       = "transport"
)

type CtxRequestIdKey struct{}
type CtxSessionIdKey struct{}
type CtxRootSpanKey struct{}

type MCPSSEContextHandler struct {
	generator uuid.UuidGenerator
}

func NewMCPSSEContextHandler(generator uuid.UuidGenerator) *MCPSSEContextHandler {
	return &MCPSSEContextHandler{
		generator: generator,
	}
}

func (h *MCPSSEContextHandler) Handle() server.SSEContextFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		sID := r.URL.Query().Get("sessionId")

		ctx = context.WithValue(ctx, CtxSessionIdKey{}, sID)

		rID := r.Header.Get(HeaderRequestID)

		if rID == "" {
			rID = h.generator.Generate()
			r.Header.Set(HeaderRequestID, rID)
		}

		ctx = context.WithValue(ctx, CtxRequestIdKey{}, rID)

		ctx, span := trace.CtxTracer(ctx).Start(
			ctx,
			"MCP",
			oteltrace.WithNewRoot(),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(
				attribute.String(SessionID, sID),
				attribute.String(RequestID, rID),
				attribute.String(Transport, "sse"),
			),
		)
		ctx = context.WithValue(ctx, CtxRootSpanKey{}, span)

		logger := log.CtxLogger(ctx).
			With().
			Str(SessionID, sID).
			Str(RequestID, rID).
			Str(Transport, "sse").
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

type MCPSSEServer struct {
	server *server.SSEServer
}

func NewMCPSSEServer(mcpServer *server.MCPServer, opts ...server.SSEOption) *MCPSSEServer {
	return &MCPSSEServer{
		server: server.NewSSEServer(mcpServer, opts...),
	}
}

func (s *MCPSSEServer) Start(ctx context.Context, addr string) error {
	log.CtxLogger(ctx).Info().Msgf("starting MCP SSE server on %s", addr)

	return s.server.Start(addr)
}

func (s *MCPSSEServer) Stop(ctx context.Context) error {
	log.CtxLogger(ctx).Info().Msg("stopping MCP SSE server")

	return s.server.Shutdown(ctx)
}
