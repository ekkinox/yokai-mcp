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
	"github.com/prometheus/client_golang/prometheus"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "mcpserver"

var MCPServerModule = fx.Module(
	ModuleName,
	fx.Provide(
		// module fixed dependencies
		ProvideMCPServerRegistry,
		ProvideMCPServer,
		ProvideMCPSSEServer,
		ProvideMCPStdioServer,
		// module overridable dependencies
		fx.Annotate(
			ProvideDefaultMCPServerHooksProvider,
			fx.As(new(yokaimcpserver.MCPServerHooksProvider)),
		),
		fx.Annotate(
			ProvideDefaultMCPServerFactory,
			fx.As(new(yokaimcpserver.MCPServerFactory)),
		),
		fx.Annotate(
			ProvideDefaultMCPSSEServerContextHandler,
			fx.As(new(sse.MCPSSEServerContextHandler)),
		),
		fx.Annotate(
			ProvideDefaultMCPSSEServerFactory,
			fx.As(new(sse.MCPSSEServerFactory)),
		),
		fx.Annotate(
			ProvideDefaultMCPStdioServerContextHandler,
			fx.As(new(stdio.MCPStdioServerContextHandler)),
		),
		fx.Annotate(
			ProvideDefaultMCPStdioServerFactory,
			fx.As(new(stdio.MCPStdioServerFactory)),
		),
		// module info
		fx.Annotate(
			NewMCPServerModuleInfo,
			fx.As(new(any)),
			fx.ResultTags(`group:"core-module-infos"`),
		),
	),
)

type ProvideDefaultMCPServerHooksProviderParams struct {
	fx.In
	Registry *prometheus.Registry
	Config   *config.Config
}

func ProvideDefaultMCPServerHooksProvider(p ProvideDefaultMCPServerHooksProviderParams) *yokaimcpserver.DefaultMCPServerHooksProvider {
	return yokaimcpserver.NewDefaultMCPServerHooksProvider(p.Registry, p.Config)
}

type ProvideDefaultMCPServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideDefaultMCPServerFactory(p ProvideDefaultMCPServerFactoryParams) *yokaimcpserver.DefaultMCPServerFactory {
	return yokaimcpserver.NewDefaultMCPServerFactory(p.Config)
}

type ProvideMCPServerRegistryParams struct {
	fx.In
	Config            *config.Config
	Tools             []yokaimcpserver.MCPServerTool             `group:"mcp-server-tools"`
	Prompts           []yokaimcpserver.MCPServerPrompt           `group:"mcp-server-prompts"`
	Resources         []yokaimcpserver.MCPServerResource         `group:"mcp-server-resources"`
	ResourceTemplates []yokaimcpserver.MCPServerResourceTemplate `group:"mcp-server-resource-templates"`
}

func ProvideMCPServerRegistry(p ProvideMCPServerRegistryParams) *yokaimcpserver.MCPServerRegistry {
	return yokaimcpserver.NewMCPServerRegistry(
		p.Config,
		p.Tools,
		p.Prompts,
		p.Resources,
		p.ResourceTemplates,
	)
}

type ProvideMCPServerParam struct {
	fx.In
	Config   *config.Config
	Provider yokaimcpserver.MCPServerHooksProvider
	Factory  yokaimcpserver.MCPServerFactory
	Registry *yokaimcpserver.MCPServerRegistry
}

func ProvideMCPServer(p ProvideMCPServerParam) *server.MCPServer {
	srv := p.Factory.Create(server.WithHooks(p.Provider.Provide()))

	p.Registry.Register(srv)

	return srv
}

type ProvideDefaultMCPSSEContextHandlerParam struct {
	fx.In
	Generator      uuid.UuidGenerator
	TracerProvider oteltrace.TracerProvider
	Logger         *log.Logger
}

func ProvideDefaultMCPSSEServerContextHandler(p ProvideDefaultMCPSSEContextHandlerParam) *sse.DefaultMCPSSEServerContextHandler {
	return sse.NewDefaultMCPSSEServerContextHandler(p.Generator, p.TracerProvider, p.Logger)
}

type ProvideDefaultMCPSSEServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideDefaultMCPSSEServerFactory(p ProvideDefaultMCPServerFactoryParams) *sse.DefaultMCPSSEServerFactory {
	return sse.NewDefaultMCPSSEServerFactory(p.Config)
}

type ProvideMCPSSEServerParam struct {
	fx.In
	LifeCycle                  fx.Lifecycle
	Context                    context.Context
	Logger                     *log.Logger
	Config                     *config.Config
	MCPServer                  *server.MCPServer
	MCPSSEServerFactory        sse.MCPSSEServerFactory
	MCPSSEServerContextHandler sse.MCPSSEServerContextHandler
}

func ProvideMCPSSEServer(p ProvideMCPSSEServerParam) *sse.MCPSSEServer {
	sseServer := p.MCPSSEServerFactory.Create(
		p.MCPServer,
		server.WithSSEContextFunc(p.MCPSSEServerContextHandler.Handle()),
	)

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

type ProvideDefaultMCPStdioContextHandlerParam struct {
	fx.In
	Generator      uuid.UuidGenerator
	TracerProvider oteltrace.TracerProvider
	Logger         *log.Logger
}

func ProvideDefaultMCPStdioServerContextHandler(p ProvideDefaultMCPStdioContextHandlerParam) *stdio.DefaultMCPStdioServerContextHandler {
	return stdio.NewDefaultMCPStdioServerContextHandler(p.Generator, p.TracerProvider, p.Logger)
}

type ProvideDefaultMCPStdioServerFactoryParams struct {
	fx.In
	Config *config.Config
}

func ProvideDefaultMCPStdioServerFactory(p ProvideDefaultMCPStdioServerFactoryParams) *stdio.DefaultMCPStdioServerFactory {
	return stdio.NewDefaultMCPStdioServerFactory(p.Config)
}

type ProvideMCPStdioServerParam struct {
	fx.In
	LifeCycle                    fx.Lifecycle
	Context                      context.Context
	Logger                       *log.Logger
	Config                       *config.Config
	MCPServer                    *server.MCPServer
	MCPStdioServerFactory        stdio.MCPStdioServerFactory
	MCPStdioServerContextHandler stdio.MCPStdioServerContextHandler
}

func ProvideMCPStdioServer(p ProvideMCPStdioServerParam) *stdio.MCPStdioServer {
	stdioServer := p.MCPStdioServerFactory.Create(
		p.MCPServer,
		server.WithStdioContextFunc(p.MCPStdioServerContextHandler.Handle()),
	)

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
