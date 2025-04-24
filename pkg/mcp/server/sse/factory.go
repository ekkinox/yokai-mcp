package sse

import (
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultAddr              = ":8082"
	DefaultBaseURL           = ""
	DefaultBasePath          = ""
	DefaultSSEEndpoint       = "/sse"
	DefaultMessageEndpoint   = "/message"
	DefaultKeepAliveInterval = 10 * time.Second
)

type MCPSSEServerFactory struct {
	config *config.Config
}

func NewMCPSSEServerFactory(config *config.Config) *MCPSSEServerFactory {
	return &MCPSSEServerFactory{
		config: config,
	}
}

func (f *MCPSSEServerFactory) Create(mcpServer *server.MCPServer, options ...server.SSEOption) *MCPSSEServer {
	addr := f.config.GetString("modules.mcp.server.transport.sse.address")
	if addr == "" {
		addr = DefaultAddr
	}

	baseURL := f.config.GetString("modules.mcp.server.transport.sse.base_url")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	basePath := f.config.GetString("modules.mcp.server.transport.sse.base_path")
	if basePath == "" {
		basePath = DefaultBasePath
	}

	sseEndpoint := f.config.GetString("modules.mcp.server.transport.sse.sse_endpoint")
	if sseEndpoint == "" {
		sseEndpoint = DefaultSSEEndpoint
	}

	messageEndpoint := f.config.GetString("modules.mcp.server.transport.sse.message_endpoint")
	if messageEndpoint == "" {
		messageEndpoint = DefaultMessageEndpoint
	}

	keepAliveInterval := DefaultKeepAliveInterval
	keepAliveIntervalConfig := f.config.GetInt("modules.mcp.server.transport.sse.keep_alive_interval")
	if keepAliveIntervalConfig != 0 {
		keepAliveInterval = time.Duration(keepAliveIntervalConfig) * time.Second
	}

	srvOptions := []server.SSEOption{
		server.WithBaseURL(baseURL),
		server.WithBasePath(basePath),
		server.WithSSEEndpoint(sseEndpoint),
		server.WithMessageEndpoint(messageEndpoint),
	}

	if f.config.GetBool("modules.mcp.server.transport.sse.keep_alive") {
		srvOptions = append(srvOptions, server.WithKeepAliveInterval(keepAliveInterval))
	}

	srvOptions = append(srvOptions, options...)

	return NewMCPSSEServer(mcpServer, addr, options...)
}
