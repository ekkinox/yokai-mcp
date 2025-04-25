package server

import (
	"reflect"
	"runtime"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServerTool interface {
	Name() string
	Options() []mcp.ToolOption
	Handle() server.ToolHandlerFunc
}

type MCPServerPrompt interface {
	Name() string
	Options() []mcp.PromptOption
	Handle() server.PromptHandlerFunc
}

type MCPServerResource interface {
	Name() string
	URI() string
	Options() []mcp.ResourceOption
	Handle() server.ResourceHandlerFunc
}

type MCPServerResourceTemplate interface {
	Name() string
	URI() string
	Options() []mcp.ResourceTemplateOption
	Handle() server.ResourceTemplateHandlerFunc
}

type MCPServerRegistryInfo struct {
	Tools             map[string]string
	Prompts           map[string]string
	Resources         map[string]string
	ResourceTemplates map[string]string
}

type MCPServerRegistry struct {
	tools             map[string]MCPServerTool
	prompts           map[string]MCPServerPrompt
	resources         map[string]MCPServerResource
	resourceTemplates map[string]MCPServerResourceTemplate
}

func NewMCPServerRegistry(
	tools []MCPServerTool,
	prompts []MCPServerPrompt,
	resources []MCPServerResource,
	resourceTemplates []MCPServerResourceTemplate,
) *MCPServerRegistry {
	toolsMap := make(map[string]MCPServerTool, len(tools))
	promptsMap := make(map[string]MCPServerPrompt, len(prompts))
	resourcesMap := make(map[string]MCPServerResource, len(resources))
	resourceTemplatesMap := make(map[string]MCPServerResourceTemplate, len(resourceTemplates))

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

	return &MCPServerRegistry{
		tools:             toolsMap,
		prompts:           promptsMap,
		resources:         resourcesMap,
		resourceTemplates: resourceTemplatesMap,
	}
}

func (r *MCPServerRegistry) Register(mcpServer *server.MCPServer) {
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

func (r *MCPServerRegistry) Info() MCPServerRegistryInfo {
	reflectInfo := func(x any) string {
		t := reflect.ValueOf(x).Type()
		if t.Kind() == reflect.Func {
			return runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
		}

		return t.String()
	}

	toolsInfo := make(map[string]string, len(r.tools))
	for _, tool := range r.tools {
		toolsInfo[tool.Name()] = reflectInfo(tool.Handle())
	}

	promptsInfo := make(map[string]string, len(r.prompts))
	for _, prompt := range r.prompts {
		promptsInfo[prompt.Name()] = reflectInfo(prompt.Handle())
	}

	resourcesInfo := make(map[string]string, len(r.resources))
	for _, resource := range r.resources {
		resourcesInfo[resource.Name()] = reflectInfo(resource.Handle())
	}

	resourceTemplatesInfo := make(map[string]string, len(r.resourceTemplates))
	for _, resourceTemplate := range r.resourceTemplates {
		resourceTemplatesInfo[resourceTemplate.Name()] = reflectInfo(resourceTemplate.Handle())
	}

	return MCPServerRegistryInfo{
		Tools:             toolsInfo,
		Prompts:           promptsInfo,
		Resources:         resourcesInfo,
		ResourceTemplates: resourceTemplatesInfo,
	}
}
