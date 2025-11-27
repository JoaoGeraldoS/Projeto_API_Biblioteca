package categories

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type CategoryRepository struct {
	exec persistence.Executer
}

func NewCategoryRepository(exec persistence.Executer) *CategoryRepository {
	return &CategoryRepository{exec: exec}
}

func (r *CategoryRepository) WithTx(exec persistence.Executer) CategoryRepositoryTx {
	return &CategoryRepository{exec: exec}
}

func (r *CategoryRepository) Create(ctx context.Context, c *Category) error {
	query := "INSERT INTO categories (name) VALUES (?)"

	resul, err := r.exec.ExecContext(ctx, query, c.Name)
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
	rows, err := r.exec.QueryContext(ctx, "SELECT * FROM categories")
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
