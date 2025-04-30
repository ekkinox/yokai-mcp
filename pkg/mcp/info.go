package mcp

import (
	"github.com/ankorstore/yokai/config"
	yokaimcpserver "github.com/ekkinox/yokai-mcp/pkg/mcp/server"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/sse"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/stdio"
)

type MCPServerModuleInfo struct {
	config      *config.Config
	registry    *yokaimcpserver.MCPServerRegistry
	sseServer   *sse.MCPSSEServer
	stdioServer *stdio.MCPStdioServer
}

func NewMCPServerModuleInfo(
	config *config.Config,
	registry *yokaimcpserver.MCPServerRegistry,
	sseServer *sse.MCPSSEServer,
	stdioServer *stdio.MCPStdioServer,
) *MCPServerModuleInfo {
	return &MCPServerModuleInfo{
		config:      config,
		registry:    registry,
		sseServer:   sseServer,
		stdioServer: stdioServer,
	}
}

// Name return the name of the module info.
func (i *MCPServerModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *MCPServerModuleInfo) Data() map[string]interface{} {
	sseServerInfo := i.sseServer.Info()
	stdioServerInfo := i.stdioServer.Info()
	mcpRegistryInfo := i.registry.Info()

	return map[string]interface{}{
		"transports": map[string]interface{}{
			"sse":   sseServerInfo,
			"stdio": stdioServerInfo,
		},
		"capabilities": map[string]interface{}{
			"tools":     mcpRegistryInfo.Capabilities.Tools,
			"prompts":   mcpRegistryInfo.Capabilities.Prompts,
			"resources": mcpRegistryInfo.Capabilities.Resources,
		},
		"registrations": map[string]interface{}{
			"tools":             mcpRegistryInfo.Registrations.Tools,
			"prompts":           mcpRegistryInfo.Registrations.Prompts,
			"resources":         mcpRegistryInfo.Registrations.Resources,
			"resourceTemplates": mcpRegistryInfo.Registrations.ResourceTemplates,
		},
	}
}
