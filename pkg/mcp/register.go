package mcp

import "go.uber.org/fx"

func AsMCPTool(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(MCPTool)),
			fx.ResultTags(`group:"mcp-tools"`),
		),
	)
}

func AsMCPPrompt(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(MCPPrompt)),
			fx.ResultTags(`group:"mcp-prompts"`),
		),
	)
}

func AsMCPResource(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(MCPResource)),
			fx.ResultTags(`group:"mcp-resources"`),
		),
	)
}

func AsMCPResourceTemplate(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(MCPResourceTemplate)),
			fx.ResultTags(`group:"mcp-resource-templates"`),
		),
	)
}
