package domain

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/trace"
)

type BookService struct {
	repository *BookRepository
}

func NewBookService(repository *BookRepository) *BookService {
	return &BookService{
		repository: repository,
	}
}

type ListBooksParams struct {
	Genre string
}

func (s *BookService) ListBooks(ctx context.Context, params ListBooksParams) ([]Book, error) {
	ctx, span := trace.CtxTracer(ctx).Start(ctx, "BookService ListBooks")
	defer span.End()

	return s.repository.Select(ctx, SelectParams{
		Genre: params.Genre,
	})
}

type CreateBookParams struct {
	Title    string
	Genre    string
	Synopsis string
}

func (s *BookService) CreateBook(ctx context.Context, params CreateBookParams) (Book, error) {
	ctx, span := trace.CtxTracer(ctx).Start(ctx, "BookService CreateBook")
	defer span.End()

	var book Book

	lastInsertID, err := s.repository.Insert(ctx, InsertParams{
		Title:    params.Title,
		Genre:    params.Genre,
		Synopsis: params.Synopsis,
	})
	if err != nil {
		return book, err
	}

	books, err := s.repository.Select(ctx, SelectParams{
		ID: lastInsertID,
	})
	if err != nil {
		return book, err
	}

	if len(books) != 1 {
		return book, fmt.Errorf("expected 1 book, got %d", len(books))
	}

	return books[0], nil
}

type DeleteBookParams struct {
	ID    int
	Genre string
}

func (s *BookService) DeleteBook(ctx context.Context, params DeleteBookParams) (int, error) {
	ctx, span := trace.CtxTracer(ctx).Start(ctx, "BookService DeleteBook")
	defer span.End()

	return s.repository.Delete(ctx, DeleteParams{
		ID:    params.ID,
		Genre: params.Genre,
	})
}
