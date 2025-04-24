package mcp

import (
	"github.com/ankorstore/yokai/config"
	yokaimcpserver "github.com/ekkinox/yokai-mcp/pkg/mcp/server"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/sse"
	"github.com/ekkinox/yokai-mcp/pkg/mcp/server/stdio"
	"github.com/mark3labs/mcp-go/server"
)

type MCPModuleInfo struct {
	config      *config.Config
	mspServer   *server.MCPServer
	mcpRegistry *yokaimcpserver.MCPServerRegistry
	sseServer   *sse.MCPSSEServer
	stdioServer *stdio.MCPStdioServer
}

func NewMCPModuleInfo(
	config *config.Config,
	mspServer *server.MCPServer,
	mcpRegistry *yokaimcpserver.MCPServerRegistry,
	sseServer *sse.MCPSSEServer,
	stdioServer *stdio.MCPStdioServer,

) *MCPModuleInfo {
	return &MCPModuleInfo{
		config:      config,
		mspServer:   mspServer,
		mcpRegistry: mcpRegistry,
		sseServer:   sseServer,
		stdioServer: stdioServer,
	}
}

// Name return the name of the module info.
func (i *MCPModuleInfo) Name() string {
	return ModuleName
}

// Data return the data of the module info.
func (i *MCPModuleInfo) Data() map[string]interface{} {
	registryInfo := i.mcpRegistry.Info()

	return map[string]interface{}{
		"server": map[string]interface{}{
			"transport": map[string]interface{}{
				"sse":   i.sseServer.Info(),
				"stdio": i.stdioServer.Info(),
			},
			"tools":             registryInfo.Tools,
			"prompts":           registryInfo.Prompts,
			"resources":         registryInfo.Resources,
			"resourceTemplates": registryInfo.ResourceTemplates,
		},
	}
}
