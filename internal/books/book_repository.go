package books

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type BookRepository struct {
	exec persistence.Executer
}

func NewBookRepository(exec persistence.Executer) *BookRepository {
	return &BookRepository{exec: exec}
}

func (r *BookRepository) WithTx(exec persistence.Executer) BookRepositoryTx {
	return &BookRepository{exec: exec}
}

func (r *BookRepository) Create(ctx context.Context, b *Books) error {
	query := "INSERT INTO books (title, description, content, author_id) VALUES (?, ?, ?, ?)"

	result, err := r.exec.ExecContext(ctx, query, b.Title, b.Description, b.Content, b.AuthorID)
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

	_, err := r.exec.ExecContext(ctx, query, book_id, cat_id)
	if err != nil {
		return err
	}

	return nil
}
