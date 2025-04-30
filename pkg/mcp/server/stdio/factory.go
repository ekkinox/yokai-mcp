package stdio

import (
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

var _ MCPStdioServerFactory = (*DefaultMCPStdioServerFactory)(nil)

type MCPStdioServerFactory interface {
	Create(mcpServer *server.MCPServer, options ...server.StdioOption) *MCPStdioServer
}

type DefaultMCPStdioServerFactory struct {
	config *config.Config
}

func NewDefaultMCPStdioServerFactory(config *config.Config) *DefaultMCPStdioServerFactory {
	return &DefaultMCPStdioServerFactory{
		config: config,
	}
}

func (f *DefaultMCPStdioServerFactory) Create(mcpServer *server.MCPServer, options ...server.StdioOption) *MCPStdioServer {
	srvConfig := MCPStdioServerConfig{
		In:  os.Stdin,
		Out: os.Stdout,
	}

	return NewMCPStdioServer(mcpServer, srvConfig, options...)
}
