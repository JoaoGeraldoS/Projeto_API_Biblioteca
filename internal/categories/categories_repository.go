package categories

import (
	"context"
	"database/sql"
	"errors"
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

		rows.Scan(&c.ID, &c.Name, &c.CreatedAT)
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoryRepository) GetById(ctx context.Context, id int64) (*Category, error) {
	var c Category

	err := r.db.QueryRowContext(ctx, `SELECT * FROM categories WHERE id = ?`, id).
		Scan(&c.ID, &c.Name, &c.CreatedAT)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(ctx context.Context, c *Category) error {
	sql := "UPDATE categories SET name = ? WHERE id = ?"

	result, err := r.db.ExecContext(ctx, sql, c.Name, c.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("erro ao atualizar categoria")
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("erro ao apagar categoria")
	}
	return nil
}
