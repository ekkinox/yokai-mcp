package mcp

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	yokaimcpserver "github.com/ekkinox/yokai-mcp/pkg/mcp/server"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/sse"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/stdio"
	"github.com/mark3labs/mcp-go/server"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "mcp"

var MCPModule = fx.Module(
	ModuleName,
	fx.Provide(
		// dependencies
		ProvideMCPServerHooksProvider,
		ProvideMCPServerFactory,
		ProvideMCPServerRegistry,
		ProvideMCPServer,
		ProvideMCPSSEServerContextHandler,
		ProvideMCPSSEServerFactory,
		ProvideMCPSSEServer,
		ProvideMCPStdioServerContextHandler,
		ProvideMCPStdioServerFactory,
		ProvideMCPStdioServer,
		// info
		fx.Annotate(
			NewMCPModuleInfo,
			fx.As(new(interface{})),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
	fx.Invoke(func(*sse.MCPSSEServer, *stdio.MCPStdioServer) {}),
)

type ProvideMCPServerHooksProviderParams struct {
	fx.In
	Config *config.Config
}

func ProvideMCPServerHooksProvider(p ProvideMCPServerHooksProviderParams) *yokaimcpserver.MCPServerHooksProvider {
	return yokaimcpserver.NewMCPServerHooksProvider(p.Config)
}

type ProvideMCPServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideMCPServerFactory(p ProvideMCPServerFactoryParams) *yokaimcpserver.MCPServerFactory {
	return yokaimcpserver.NewMCPServerFactory(p.Config)
}

type ProvideMCPServerRegistryParams struct {
	fx.In
	Tools             []yokaimcpserver.MCPServerTool             `group:"mcp-server-tools"`
	Prompts           []yokaimcpserver.MCPServerPrompt           `group:"mcp-server-prompts"`
	Resources         []yokaimcpserver.MCPServerResource         `group:"mcp-server-resources"`
	ResourceTemplates []yokaimcpserver.MCPServerResourceTemplate `group:"mcp-server-resource-templates"`
}

func ProvideMCPServerRegistry(p ProvideMCPServerRegistryParams) *yokaimcpserver.MCPServerRegistry {
	return yokaimcpserver.NewMCPServerRegistry(p.Tools, p.Prompts, p.Resources, p.ResourceTemplates)
}

type ProvideMCPServerParam struct {
	fx.In
	Config   *config.Config
	Provider *yokaimcpserver.MCPServerHooksProvider
	Factory  *yokaimcpserver.MCPServerFactory
	Registry *yokaimcpserver.MCPServerRegistry
}

func ProvideMCPServer(p ProvideMCPServerParam) *server.MCPServer {
	srv := p.Factory.Create(server.WithHooks(p.Provider.Provide()))

	p.Registry.Register(srv)

	return srv
}

type ProvideMCPSSEContextHandlerParam struct {
	fx.In
	Generator      uuid.UuidGenerator
	TracerProvider oteltrace.TracerProvider
	Logger         *log.Logger
}

func ProvideMCPSSEServerContextHandler(p ProvideMCPSSEContextHandlerParam) *sse.MCPSSEServerContextHandler {
	return sse.NewMCPSSEServerContextHandler(p.Generator, p.TracerProvider, p.Logger)
}

type ProvideMCPSSEServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideMCPSSEServerFactory(p ProvideMCPServerFactoryParams) *sse.MCPSSEServerFactory {
	return sse.NewMCPSSEServerFactory(p.Config)
}

type ProvideMCPSSEServerParam struct {
	fx.In
	LifeCycle                  fx.Lifecycle
	Context                    context.Context
	Logger                     *log.Logger
	Config                     *config.Config
	MCPServer                  *server.MCPServer
	MCPSSEServerFactory        *sse.MCPSSEServerFactory
	MCPSSEServerContextHandler *sse.MCPSSEServerContextHandler
}

func ProvideMCPSSEServer(p ProvideMCPSSEServerParam) *sse.MCPSSEServer {
	sseServer := p.MCPSSEServerFactory.Create(p.MCPServer, server.WithSSEContextFunc(p.MCPSSEServerContextHandler.Handle()))

	if p.Config.GetBool("modules.mcp.server.transport.sse.expose") {
		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go sseServer.Start(p.Context)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				return sseServer.Stop(ctx)
			},
		})
	}

	return sseServer
}

type ProvideMCPStdioContextHandlerParam struct {
	fx.In
	Generator      uuid.UuidGenerator
	TracerProvider oteltrace.TracerProvider
	Logger         *log.Logger
}

func ProvideMCPStdioServerContextHandler(p ProvideMCPStdioContextHandlerParam) *stdio.MCPStdioServerContextHandler {
	return stdio.NewMCPStdioServerContextHandler(p.Generator, p.TracerProvider, p.Logger)
}

type ProvideMCPStdioServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideMCPStdioServerFactory(p ProvideMCPStdioServerFactoryParams) *stdio.MCPStdioServerFactory {
	return stdio.NewMCPStdioServerFactory(p.Config)
}

type ProvideMCPStdioServerParam struct {
	fx.In
	LifeCycle                    fx.Lifecycle
	Context                      context.Context
	Logger                       *log.Logger
	Config                       *config.Config
	MCPServer                    *server.MCPServer
	MCPStdioServerFactory        *stdio.MCPStdioServerFactory
	MCPStdioServerContextHandler *stdio.MCPStdioServerContextHandler
}

func ProvideMCPStdioServer(p ProvideMCPStdioServerParam) *stdio.MCPStdioServer {
	stdioServer := p.MCPStdioServerFactory.Create(p.MCPServer, server.WithStdioContextFunc(p.MCPStdioServerContextHandler.Handle()))

	if p.Config.GetBool("modules.mcp.server.transport.stdio.expose") {
		p.LifeCycle.Append(fx.Hook{
			OnStart: func(context.Context) error {
				go stdioServer.Start(p.Context)

				return nil
			},
		})
	}

	return stdioServer
}
