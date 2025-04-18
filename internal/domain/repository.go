package domain

import (
	"context"
	"database/sql"

	"github.com/huandu/go-sqlbuilder"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}

type SelectParams struct {
	ID    int
	Genre string
}

func (r *BookRepository) Select(ctx context.Context, params SelectParams) ([]Book, error) {
	var books []Book

	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(
		"id",
		"title",
		"genre",
		"synopsis",
	)
	sb.From("books")
	sb.OrderBy("id ASC")

	if params.ID != 0 {
		sb.Where(sb.Equal("id", params.ID))
	}

	if params.Genre != "" {
		sb.Where(sb.Equal("genre", params.Genre))
	}

	query, args := sb.Build()

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book

		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Genre,
			&book.Synopsis,
		); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, rows.Err()
}

type InsertParams struct {
	Title    string
	Genre    string
	Synopsis string
}

func (r *BookRepository) Insert(ctx context.Context, params InsertParams) (int, error) {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("books")
	ib.Cols(
		"title",
		"genre",
		"synopsis",
	)
	ib.Values(
		params.Title,
		params.Genre,
		params.Synopsis,
	)
	query, args := ib.Build()

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()

	return int(lastInsertId), err
}

type DeleteParams struct {
	ID    int
	Genre string
}

func (r *BookRepository) Delete(ctx context.Context, params DeleteParams) (int, error) {
	db := sqlbuilder.NewDeleteBuilder()
	db.DeleteFrom("books")

	if params.ID != 0 {
		db.Where(db.Equal("id", params.ID))
	}

	if params.Genre != "" {
		db.Where(db.Equal("genre", params.Genre))
	}

	query, args := db.Build()

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()

	return int(rowsAffected), err
}
