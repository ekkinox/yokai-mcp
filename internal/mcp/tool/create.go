package tool

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ekkinox/yokai-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type CreateBookTool struct {
	service *domain.BookService
}

func NewCreateBookTool(service *domain.BookService) *CreateBookTool {
	return &CreateBookTool{
		service: service,
	}
}

func (t *CreateBookTool) Name() string {
	return "create-book"
}

func (t *CreateBookTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("To create a new book."),
		mcp.WithString(
			"title",
			mcp.Required(),
			mcp.Description("Title of the book."),
		),
		mcp.WithString(
			"genre",
			mcp.Required(),
			mcp.Description("Genre of the book."),
			mcp.Enum("science-fiction", "horror", "romance", "fantasy"),
		),
		mcp.WithString(
			"synopsis",
			mcp.Required(),
			mcp.Description("Synopsis of the book."),
		),
	}
}

func (t *CreateBookTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		titleParams, ok := request.Params.Arguments["title"].(string)
		if !ok {
			return nil, errors.New("title must be a string")
		}

		genreParam, ok := request.Params.Arguments["genre"].(string)
		if !ok {
			return nil, errors.New("genre must be a string")
		}

		synopsisParam, ok := request.Params.Arguments["synopsis"].(string)
		if !ok {
			return nil, errors.New("synopsis must be a string")
		}

		book, err := t.service.CreateBook(ctx, domain.CreateBookParams{
			Title:    titleParams,
			Genre:    genreParam,
			Synopsis: synopsisParam,
		})
		if err != nil {
			return nil, err
		}

		jsonBook, err := json.Marshal(book)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(jsonBook)), nil
	}
}
