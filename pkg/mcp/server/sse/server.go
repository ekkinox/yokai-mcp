package sse

import (
	"context"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

type MCPSSEServerConfig struct {
	Address           string
	BaseURL           string
	BasePath          string
	SSEEndpoint       string
	MessageEndpoint   string
	KeepAlive         bool
	KeepAliveInterval time.Duration
}

type MCPSSEServer struct {
	server  *server.SSEServer
	config  MCPSSEServerConfig
	running bool
}

func NewMCPSSEServer(mcpServer *server.MCPServer, config MCPSSEServerConfig, opts ...server.SSEOption) *MCPSSEServer {
	return &MCPSSEServer{
		server: server.NewSSEServer(mcpServer, opts...),
		config: config,
	}
}

func (s *MCPSSEServer) Server() *server.SSEServer {
	return s.server
}

func (s *MCPSSEServer) Config() MCPSSEServerConfig {
	return s.config
}

func (s *MCPSSEServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msgf("starting MCP SSE server on %s", s.config.Address)

	s.running = true

	err := s.server.Start(s.config.Address)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP SSE server")

		s.running = false
	}

	return err
}

func (s *MCPSSEServer) Stop(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("stopping MCP SSE server")

	err := s.server.Shutdown(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to stop MCP SSE server")
	}

	s.running = false

	return err
}

func (s *MCPSSEServer) Running() bool {
	return s.running
}

func (s *MCPSSEServer) Info() map[string]any {
	return map[string]any{
		"config": map[string]any{
			"address":             s.config.Address,
			"base_url":            s.config.BaseURL,
			"base_path":           s.config.BasePath,
			"sse_endpoint":        s.config.SSEEndpoint,
			"message_endpoint":    s.config.MessageEndpoint,
			"keep_alive":          s.config.KeepAlive,
			"keep_alive_interval": s.config.KeepAliveInterval.Seconds(),
		},
		"status": map[string]any{
			"running": s.running,
		},
	}
}
