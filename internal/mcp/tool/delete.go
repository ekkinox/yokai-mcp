package tool

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ekkinox/yokai-mcp/internal/domain"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type DeleteBookTool struct {
	service *domain.BookService
}

func NewDeleteBookTool(service *domain.BookService) *DeleteBookTool {
	return &DeleteBookTool{
		service: service,
	}
}

func (t *DeleteBookTool) Name() string {
	return "delete-book"
}

func (t *DeleteBookTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("To delete one or several existing books."),
		mcp.WithString(
			"id",
			mcp.DefaultString(""),
			mcp.Description("Optional ID of the book to delete. Empty value means no book selection by id."),
		),
		mcp.WithString(
			"genre",
			mcp.DefaultString(""),
			mcp.Description("Optional genre of the book. Empty value means bo books selection by genre."),
			mcp.Enum("", "science-fiction", "horror", "romance", "fantasy"),
		),
	}
}

func (t *DeleteBookTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := 0
		idParam, ok := request.Params.Arguments["id"]
		if ok {
			var err error

			id, err = strconv.Atoi(idParam.(string))
			if err != nil {
				return nil, err
			}
		}

		genre := ""
		genreParam, ok := request.Params.Arguments["genre"]
		if ok {
			genre = genreParam.(string)
		}

		rowsAffected, err := t.service.DeleteBook(ctx, domain.DeleteBookParams{
			ID:    id,
			Genre: genre,
		})
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(fmt.Sprintf("%d books were deleted", rowsAffected)), nil
	}
}
