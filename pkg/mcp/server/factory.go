package server

import (
	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/server"
)

const (
	DefaultServerName    = "yokai-mcp"
	DefaultServerVersion = "1.0.0"
)

type MCPServerFactory struct {
	config *config.Config
}

func NewMCPServerFactory(config *config.Config) *MCPServerFactory {
	return &MCPServerFactory{
		config: config,
	}
}

func (f *MCPServerFactory) Create(options ...server.ServerOption) *server.MCPServer {
	name := f.config.GetString("modules.mcp.server.name")
	if name == "" {
		name = DefaultServerName
	}

	version := f.config.GetString("modules.mcp.server.version")
	if version == "" {
		version = DefaultServerVersion
	}

	srvOptions := []server.ServerOption{
		server.WithRecovery(),
	}

	if f.config.GetBool("modules.mcp.server.capabilities.resource") {
		srvOptions = append(srvOptions, server.WithResourceCapabilities(true, true))
	}

	if f.config.GetBool("modules.mcp.server.capabilities.prompt") {
		srvOptions = append(srvOptions, server.WithPromptCapabilities(true))
	}

	if f.config.GetBool("modules.mcp.server.capabilities.prompt") {
		srvOptions = append(srvOptions, server.WithToolCapabilities(true))
	}

	if f.config.GetBool("modules.mcp.server.capabilities.logging") {
		srvOptions = append(srvOptions, server.WithLogging())
	}

	instructions := f.config.GetString("modules.mcp.server.instructions")
	if instructions != "" {
		srvOptions = append(srvOptions, server.WithInstructions(instructions))
	}

	srvOptions = append(srvOptions, options...)

	return server.NewMCPServer(name, version, srvOptions...)
}
