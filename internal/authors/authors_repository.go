package authors

import (
	"context"
	"database/sql"
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
