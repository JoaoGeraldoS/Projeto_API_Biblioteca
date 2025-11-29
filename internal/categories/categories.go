package categories

import (
	"context"
	"errors"
)

type Category struct {
	ID        int64
	Name      string
	CreatedAT string
}

type CategoryCreator interface {
	Create(ctx context.Context, c *Category) error
}

type CategoryRead interface {
	GetAll(ctx context.Context) ([]Category, error)
}

type ICategoryRepository interface {
	CategoryCreator
	CategoryRead
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("nome nao pode estar vazio")
	}

	return nil
}
