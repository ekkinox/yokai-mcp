package stdio

import (
	"context"
	"io"

	"github.com/ankorstore/yokai/log"
	"github.com/mark3labs/mcp-go/server"
)

type MCPStdioServer struct {
	server  *server.StdioServer
	in      io.Reader
	out     io.Writer
	running bool
}

func NewMCPStdioServer(mcpServer *server.MCPServer, in io.Reader, out io.Writer, opts ...server.StdioOption) *MCPStdioServer {
	stdioServer := server.NewStdioServer(mcpServer)

	for _, opt := range opts {
		opt(stdioServer)
	}

	return &MCPStdioServer{
		server: stdioServer,
		in:     in,
		out:    out,
	}
}

func (s *MCPStdioServer) Server() *server.StdioServer {
	return s.server
}

func (s *MCPStdioServer) In() io.Reader {
	return s.in
}

func (s *MCPStdioServer) Out() io.Writer {
	return s.out
}

func (s *MCPStdioServer) Start(ctx context.Context) error {
	logger := log.CtxLogger(ctx)

	logger.Info().Msg("starting MCP Stdio server")

	s.running = true

	err := s.server.Listen(ctx, s.in, s.out)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to start MCP Stdio server")

		s.running = false
	}

	return err
}

func (s *MCPStdioServer) Info() map[string]any {
	return map[string]any{
		"running": s.running,
	}
}
