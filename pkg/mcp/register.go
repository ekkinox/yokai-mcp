package mcp

import (
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server"
	"go.uber.org/fx"
)

func AsMCPServerTool(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerTool)),
			fx.ResultTags(`group:"mcp-server-tools"`),
		),
	)
}

func AsMCPServerTools(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerTool(constructor))
	}

	return fx.Options(options...)
}

func AsMCPServerPrompt(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerPrompt)),
			fx.ResultTags(`group:"mcp-server-prompts"`),
		),
	)
}

func AsMCPServerPrompts(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerPrompt(constructor))
	}

	return fx.Options(options...)
}

func AsMCPServerResource(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerResource)),
			fx.ResultTags(`group:"mcp-server-resources"`),
		),
	)
}

func AsMCPServerResources(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerResource(constructor))
	}

	return fx.Options(options...)
}

func AsMCPServerResourceTemplate(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(server.MCPServerResourceTemplate)),
			fx.ResultTags(`group:"mcp-server-resource-templates"`),
		),
	)
}

func AsMCPServerResourceTemplates(constructors ...any) fx.Option {
	options := []fx.Option{}

	for _, constructor := range constructors {
		options = append(options, AsMCPServerResourceTemplate(constructor))
	}

	return fx.Options(options...)
}
