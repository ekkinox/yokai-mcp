package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPTool interface {
	Name() string
	Options() []mcp.ToolOption
	Handle() server.ToolHandlerFunc
}

type MCPPrompt interface {
	Name() string
	Options() []mcp.PromptOption
	Handle() server.PromptHandlerFunc
}

type MCPResource interface {
	Name() string
	URI() string
	Options() []mcp.ResourceOption
	Handle() server.ResourceHandlerFunc
}

type MCPResourceTemplate interface {
	Name() string
	URI() string
	Options() []mcp.ResourceTemplateOption
	Handle() server.ResourceTemplateHandlerFunc
}

type MCPRegistry struct {
	tools             map[string]MCPTool
	prompts           map[string]MCPPrompt
	resources         map[string]MCPResource
	resourceTemplates map[string]MCPResourceTemplate
}

func NewMCPRegistry(tools []MCPTool, prompts []MCPPrompt, resources []MCPResource, resourceTemplates []MCPResourceTemplate) *MCPRegistry {
	toolsMap := make(map[string]MCPTool, len(tools))
	promptsMap := make(map[string]MCPPrompt, len(prompts))
	resourcesMap := make(map[string]MCPResource, len(resources))
	resourceTemplatesMap := make(map[string]MCPResourceTemplate, len(resourceTemplates))

	for _, tool := range tools {
		toolsMap[tool.Name()] = tool
	}

	for _, prompt := range prompts {
		promptsMap[prompt.Name()] = prompt
	}

	for _, resource := range resources {
		resourcesMap[resource.Name()] = resource
	}

	for _, resourceTemplate := range resourceTemplates {
		resourceTemplatesMap[resourceTemplate.Name()] = resourceTemplate
	}

	return &MCPRegistry{
		tools:             toolsMap,
		prompts:           promptsMap,
		resources:         resourcesMap,
		resourceTemplates: resourceTemplatesMap,
	}
}

func (r *MCPRegistry) Register(mcpServer *server.MCPServer) {
	for _, tool := range r.tools {
		mcpServer.AddTool(
			mcp.NewTool(tool.Name(), tool.Options()...),
			tool.Handle(),
		)
	}

	for _, prompt := range r.prompts {
		mcpServer.AddPrompt(
			mcp.NewPrompt(prompt.Name(), prompt.Options()...),
			prompt.Handle(),
		)
	}

	for _, resource := range r.resources {
		mcpServer.AddResource(
			mcp.NewResource(resource.URI(), resource.Name(), resource.Options()...),
			resource.Handle(),
		)
	}

	for _, resourceTemplate := range r.resourceTemplates {
		mcpServer.AddResourceTemplate(
			mcp.NewResourceTemplate(resourceTemplate.URI(), resourceTemplate.Name(), resourceTemplate.Options()...),
			resourceTemplate.Handle(),
		)
	}
}
