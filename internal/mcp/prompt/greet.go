package prompt

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type GreetPrompt struct {
	config *config.Config
}

func NewGreetPrompt(config *config.Config) *GreetPrompt {
	return &GreetPrompt{
		config: config,
	}
}

func (p *GreetPrompt) Name() string {
	return "greet"
}

func (p *GreetPrompt) Options() []mcp.PromptOption {
	return []mcp.PromptOption{
		mcp.WithPromptDescription("A books assistant greeting prompt"),
		mcp.WithArgument(
			"name",
			mcp.ArgumentDescription("Name of books owner to greet"),
		),
	}
}

func (p *GreetPrompt) Handle() server.PromptHandlerFunc {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := request.Params.Arguments["name"]
		if name == "" {
			name = p.config.GetString("config.books.default_owner")
		}

		return mcp.NewGetPromptResult(
			"A books assistant greeting",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("Hello, %s! I am your books assistant. How can I help you today?", name)),
				),
			},
		), nil
	}
}
