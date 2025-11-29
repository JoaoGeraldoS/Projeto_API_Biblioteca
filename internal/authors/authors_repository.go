package authors

import (
	"context"
	"database/sql"
	"errors"
)

type AuthorsRepository struct {
	db *sql.DB
}

func NewAuthorsRepository(db *sql.DB) *AuthorsRepository {
	return &AuthorsRepository{db: db}
}

func (r *AuthorsRepository) Create(ctx context.Context, a *Authors) error {
	query := "INSERT INTO authors (name, description) VALUES (?, ?)"
	result, err := r.db.ExecContext(ctx, query, a.Name, a.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = id
	return nil
}

func (r *AuthorsRepository) GetAll(ctx context.Context) ([]Authors, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM authors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authorsAll []Authors

	for rows.Next() {
		var a Authors
		if err := rows.Scan(&a.ID, &a.Name, &a.Description); err != nil {
			return nil, err
		}
		authorsAll = append(authorsAll, a)
	}
	return authorsAll, nil
}

func (r *AuthorsRepository) GetByID(ctx context.Context, id int64) (*Authors, error) {
	var author Authors

	err := r.db.QueryRowContext(ctx, "SELECT * FROM authors WHERE id = ?", id).
		Scan(&author.ID, &author.Name, &author.Description)
	if err != nil {
		return nil, errors.New("erro ao realizar busca")
	}
	return &author, nil
}

func (r *AuthorsRepository) Update(ctx context.Context, author *Authors) error {
	result, err := r.db.ExecContext(ctx, "UPDATE authors SET name = ?, description = ? WHERE id = ?",
		author.Name, author.Description, author.ID)
	if err != nil {
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("erro ao atualizar autor")
	}
	return nil
}

func (r *AuthorsRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM authors WHERE id = ?", id)
	if err != nil {
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return errors.New("erro ao deletar autor")
	}
	return nil
}
