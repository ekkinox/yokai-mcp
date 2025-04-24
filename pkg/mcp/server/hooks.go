package server

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/sse"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type MCPServerHooksProvider struct {
	config *config.Config
}

func NewMCPServerHooksProvider(config *config.Config) *MCPServerHooksProvider {
	return &MCPServerHooksProvider{
		config: config,
	}
}

func (p *MCPServerHooksProvider) Provide() *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(ctx context.Context, session server.ClientSession) {
		log.CtxLogger(ctx).Info().Str(sse.SessionID, session.SessionID()).Msg("MCP session registered")
	})

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		if span, ok := ctx.Value(sse.CtxRootSpanKey{}).(oteltrace.Span); ok {
			rwSpan, ok := span.(otelsdktrace.ReadWriteSpan)
			if ok {
				span.SetName(fmt.Sprintf("%s %s", rwSpan.Name(), string(method)))
			}

			span.SetAttributes(attribute.String("method", string(method)))
		}

		log.CtxLogger(ctx).Info().Msgf("MCP call start: %s, %#+v, %#+v\n", method, id, message)
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		if span, ok := ctx.Value(sse.CtxRootSpanKey{}).(oteltrace.Span); ok {
			span.SetStatus(codes.Ok, "MCP call success")
			span.End()
		}

		log.CtxLogger(ctx).Info().Msgf("MCP call success: %s, %#+v, %#+v, %#+v\n", method, id, message, result)
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		if span, ok := ctx.Value(sse.CtxRootSpanKey{}).(oteltrace.Span); ok {
			span.RecordError(err)
			span.SetStatus(codes.Error, "MCP call error")
			span.End()
		}

		log.CtxLogger(ctx).Info().Msgf("MCP call error: %s, %#+v, %#+v, %#+v\n", method, id, message, err)
	})

	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		if span, ok := ctx.Value(sse.CtxRootSpanKey{}).(oteltrace.Span); ok {
			rwSpan, ok := span.(otelsdktrace.ReadWriteSpan)
			if ok {
				span.SetName(fmt.Sprintf("%s(%s)", rwSpan.Name(), message.Params.Name))
			}

			span.SetAttributes(attribute.String("tool", message.Params.Name))
		}

		log.CtxLogger(ctx).Info().Msgf("beforeCallTool: %v, %#+v\n", id, message)
	})

	return hooks
}
