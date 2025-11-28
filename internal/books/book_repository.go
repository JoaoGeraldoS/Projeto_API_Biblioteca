package books

import (
	"context"
	"database/sql"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(ctx context.Context, b *Books) error {
	query := "INSERT INTO books (title, description, content, author_id) VALUES (?, ?, ?, ?)"

	result, err := r.db.ExecContext(ctx, query, b.Title, b.Description, b.Content, b.AuthorID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	b.ID = id
	return nil
}

func (r *BookRepository) RelationBookCategory(ctx context.Context, book_id, cat_id int64) error {
	query := "INSERT INTO book_category (book_id, category_id) VALUES (?, ?)"

	_, err := r.db.ExecContext(ctx, query, book_id, cat_id)
	if err != nil {
		return err
	}

	return nil
}
