package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ekkinox/yokai-mcp/internal/domain"
	"github.com/ekkinox/yokai-mcp/internal/http/handler"
	"github.com/ekkinox/yokai-mcp/internal/mcp/prompt"
	"github.com/ekkinox/yokai-mcp/internal/mcp/resource"
	"github.com/ekkinox/yokai-mcp/internal/mcp/tool"
	"github.com/ekkinox/yokai-mcp/pkg/mcp"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// domain
		fx.Provide(
			domain.NewBookRepository,
			domain.NewBookService,
		),
		// http
		fxhttpserver.AsHandler("GET", "/books", handler.NewListBooksHandler),
		// mcp prompts
		mcp.AsMCPPrompt(prompt.NewGreetPrompt),
		// mcp tools
		mcp.AsMCPTool(tool.NewListBooksTool),
		mcp.AsMCPTool(tool.NewCreateBookTool),
		mcp.AsMCPTool(tool.NewDeleteBookTool),
		// mcp resource
		mcp.AsMCPResourceTemplate(resource.NewOpenLibraryResource),
	)
}
