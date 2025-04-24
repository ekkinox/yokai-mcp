package stdio

import (
	"os"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

type MCPStdioServerFactory struct {
	config *config.Config
}

func NewMCPStdioServerFactory(config *config.Config) *MCPStdioServerFactory {
	return &MCPStdioServerFactory{
		config: config,
	}
}

func (f *MCPStdioServerFactory) Create(mcpServer *server.MCPServer, options ...server.StdioOption) *MCPStdioServer {
	return NewMCPStdioServer(mcpServer, os.Stdin, os.Stdout, options...)
}
