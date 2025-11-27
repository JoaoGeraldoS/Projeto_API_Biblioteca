package categories

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type Category struct {
	ID        int64
	Name      string
	CreatedAT string
}

type CategoryCreator interface {
	Create(ctx context.Context, c *Category) error
	WithTx(exec persistence.Executer) CategoryRepositoryTx
}

type CategoryRead interface {
	GetAll(ctx context.Context) ([]Category, error)
}

type CategoryRepositoryTx interface {
	CategoryCreator
	CategoryRead
}
