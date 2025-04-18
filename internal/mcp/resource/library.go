package resource

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type OpenLibraryResource struct {
	client *http.Client
}

func NewOpenLibraryResource(client *http.Client) *OpenLibraryResource {
	return &OpenLibraryResource{
		client: client,
	}
}

func (p *OpenLibraryResource) Name() string {
	return "search-book"
}

func (p *OpenLibraryResource) URI() string {
	return "books://{title}"
}

func (p *OpenLibraryResource) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("Books information search on the openlibrary.org API."),
		mcp.WithTemplateMIMEType("application/json"),
	}
}

func (p *OpenLibraryResource) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     "some fake search",
			},
		}, nil
	}
}
