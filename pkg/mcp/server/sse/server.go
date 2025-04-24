package sse

import (
	"context"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

type MCPSSEServer struct {
	server  *server.SSEServer
	address string
	running bool
}

func NewMCPSSEServer(mcpServer *server.MCPServer, address string, opts ...server.SSEOption) *MCPSSEServer {
	return &MCPSSEServer{
		server:  server.NewSSEServer(mcpServer, opts...),
		address: address,
	}
}

func (s *MCPSSEServer) Server() *server.SSEServer {
	return s.server
}

func (s *MCPSSEServer) Address() string {
	return s.address
}

func (s *MCPSSEServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msgf("starting MCP SSE server on %s", s.address)

	s.running = true

	err := s.server.Start(s.address)
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

func (s *MCPSSEServer) Info() map[string]any {
	return map[string]any{
		"address": s.address,
		"running": s.running,
	}
}
