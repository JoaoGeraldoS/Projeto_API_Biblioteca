package authors

import (
	"context"
	"errors"
)

type Authors struct {
	ID          int64
	Name        string
	Description string
}

type AuthorsCreator interface {
	Create(ctx context.Context, a *Authors) error
	Update(ctx context.Context, author *Authors) error
	Delete(ctx context.Context, id int64) error
}

type AuthorsRead interface {
	GetAll(ctx context.Context) ([]Authors, error)
	GetByID(ctx context.Context, id int64) (*Authors, error)
}

type IAuthorRepository interface {
	AuthorsCreator
	AuthorsRead
}

func (a *Authors) Validate() error {
	if a.Name == "" {
		return errors.New("Nome invalido!")
	}

	if a.Description == "" {
		return errors.New("Descricao invalida!")
	}

	return nil
}
