package resource

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ankorstore/yokai/config"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type WeatherResource struct {
	config *config.Config
	client *http.Client
}

func NewWeatherResource(config *config.Config, client *http.Client) *WeatherResource {
	return &WeatherResource{
		config: config,
		client: client,
	}
}

func (r *WeatherResource) Name() string {
	return "weather"
}

func (r *WeatherResource) URI() string {
	return "weather://{city}"
}

func (r *WeatherResource) Options() []mcp.ResourceTemplateOption {
	return []mcp.ResourceTemplateOption{
		mcp.WithTemplateDescription("Search weather information for a city on https://wttr.in/"),
		mcp.WithTemplateMIMEType("text/plain"),
	}
}

func (r *WeatherResource) Handle() server.ResourceTemplateHandlerFunc {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		city := request.Params.Arguments["city"]

		url := fmt.Sprintf("%s/%s?format=3", r.config.GetString("config.weather.host"), city)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot read body")
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/plain",
				Text:     string(body),
			},
		}, nil
	}
}
