package handler

import (
	"net/http"

	"github.com/ekkinox/yokai-mcp/internal/domain"
	"github.com/labstack/echo/v4"
)

type ListBooksHandler struct {
	service *domain.BookService
}

func NewListBooksHandler(service *domain.BookService) *ListBooksHandler {
	return &ListBooksHandler{
		service: service,
	}
}

// Handle handles HTTP requests.
func (h *ListBooksHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		books, err := h.service.ListBooks(c.Request().Context(), domain.ListBooksParams{
			Genre: c.QueryParam("genre"),
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, books)
	}
}
