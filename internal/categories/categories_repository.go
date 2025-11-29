package categories

import (
	"context"
	"database/sql"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *Category) error {
	query := "INSERT INTO categories (name) VALUES (?)"

	resul, err := r.db.ExecContext(ctx, query, c.Name)
	if err != nil {
		return err
	}

	id, err := resul.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var c Category

		rows.Scan(&c.ID, &c.Name)
		categories = append(categories, c)
	}

	return categories, nil
}
