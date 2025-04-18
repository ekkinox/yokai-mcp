package mcp

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/fx"
)

const ModuleName = "mcp"

var MCPModule = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideMCPRegistry,
		ProvideMCPServer,
		ProvideMCPSSEContextHandler,
		ProvideMCPSSEServer,
	),
	fx.Invoke(func(s *MCPSSEServer) {}),
)

type ProvideMCPRegistryParams struct {
	fx.In
	Tools             []MCPTool             `group:"mcp-tools"`
	Prompts           []MCPPrompt           `group:"mcp-prompts"`
	Resources         []MCPResource         `group:"mcp-resources"`
	ResourceTemplates []MCPResourceTemplate `group:"mcp-resource-templates"`
}

func ProvideMCPRegistry(p ProvideMCPRegistryParams) *MCPRegistry {
	return NewMCPRegistry(p.Tools, p.Prompts, p.Resources, p.ResourceTemplates)
}

type ProvideMCPServerParam struct {
	fx.In
	Config   *config.Config
	Registry *MCPRegistry
}

func ProvideMCPServer(p ProvideMCPServerParam) *server.MCPServer {
	mcpServer := server.NewMCPServer(
		p.Config.GetString("modules.mcp.name"),
		p.Config.GetString("modules.mcp.version"),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(MCPHooks()),
	)

	p.Registry.Register(mcpServer)

	return mcpServer
}

type ProvideMCPSSEContextHandlerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Generator uuid.UuidGenerator
}

func ProvideMCPSSEContextHandler(p ProvideMCPSSEContextHandlerParam) *MCPSSEContextHandler {
	return NewMCPSSEContextHandler(p.Generator)
}

type ProvideMCPSSEServerParam struct {
	fx.In
	LifeCycle         fx.Lifecycle
	Context           context.Context
	Logger            *log.Logger
	Config            *config.Config
	MCPServer         *server.MCPServer
	SSEContextHandler *MCPSSEContextHandler
}

func ProvideMCPSSEServer(p ProvideMCPSSEServerParam) *MCPSSEServer {
	sseServer := NewMCPSSEServer(
		p.MCPServer,
		server.WithSSEContextFunc(p.SSEContextHandler.Handle()),
	)

	p.LifeCycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			addr := p.Config.GetString("modules.mcp.address")

			p.Logger.Info().Msgf("starting MCP SSE server on %s", addr)

			go sseServer.Start(p.Context, p.Config.GetString("modules.mcp.address"))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info().Msg("stopping MCP SSE server")

			return sseServer.Stop(ctx)
		},
	})

	return sseServer
}
