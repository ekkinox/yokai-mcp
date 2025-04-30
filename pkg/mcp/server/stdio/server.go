package stdio

import (
	"context"
	"io"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

type MCPStdioServerConfig struct {
	In  io.Reader
	Out io.Writer
}

type MCPStdioServer struct {
	server  *server.StdioServer
	config  MCPStdioServerConfig
	running bool
}

func NewMCPStdioServer(mcpServer *server.MCPServer, config MCPStdioServerConfig, opts ...server.StdioOption) *MCPStdioServer {
	stdioServer := server.NewStdioServer(mcpServer)

	for _, opt := range opts {
		opt(stdioServer)
	}

	return &MCPStdioServer{
		server: stdioServer,
		config: config,
	}
}

func (s *MCPStdioServer) Server() *server.StdioServer {
	return s.server
}

func (s *MCPStdioServer) Config() MCPStdioServerConfig {
	return s.config
}

func (s *MCPStdioServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("starting MCP Stdio server")

	s.running = true

	err := s.server.Listen(ctx, s.config.In, s.config.Out)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP Stdio server")

		s.running = false
	}

	return err
}

func (s *MCPStdioServer) Running() bool {
	return s.running
}

func (s *MCPStdioServer) Info() map[string]any {
	return map[string]any{
		"status": map[string]any{
			"running": s.running,
		},
	}
}
