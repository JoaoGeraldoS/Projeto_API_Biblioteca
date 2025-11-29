package books

import (
	"context"
	"database/sql"
	"errors"
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

func (r *BookRepository) Update(ctx context.Context, b *Books) error {
	query := "UPDATE books SET title = ?, description = ?, content = ?, author_id = ? WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, b.Title, b.Description, b.Content, b.AuthorID, b.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("erro ao atualizar livro")
	}
	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("erro ao deletar livro")
	}
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
