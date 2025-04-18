package tool

import (
	"context"
	"encoding/json"

	"github.com/ankorstore/yokai/log"
	"github.com/ekkinox/yokai-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ListBooksTool struct {
	service *domain.BookService
}

func NewListBooksTool(service *domain.BookService) *ListBooksTool {
	return &ListBooksTool{
		service: service,
	}
}

func (t *ListBooksTool) Name() string {
	return "list-books"
}

func (t *ListBooksTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("To list one or several existing books."),
		mcp.WithString(
			"genre",
			mcp.DefaultString(""),
			mcp.Description("Optional genre of the books to list. Empty value means all genres."),
			mcp.Enum("", "science-fiction", "horror", "romance", "fantasy"),
		),
	}
}

func (t *ListBooksTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.CtxLogger(ctx).Info().Msg("some logs from the list tool")

		genre := ""
		genreParam, ok := request.Params.Arguments["genre"]
		if ok {
			genre = genreParam.(string)
		}

		books, err := t.service.ListBooks(ctx, domain.ListBooksParams{
			Genre: genre,
		})
		if err != nil {
			return nil, err
		}

		jsonBooks, err := json.Marshal(books)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(jsonBooks)), nil
	}
}
