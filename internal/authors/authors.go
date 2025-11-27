package authors

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type Authors struct {
	ID          int64
	Name        string
	Description string
}

type AuthorsCreator interface {
	Create(ctx context.Context, a *Authors) error
	WithTx(exec persistence.Executer) AuthorRepositoryTx
}

type AuthorsRead interface {
	GetAll(ctx context.Context) ([]Authors, error)
}

type AuthorRepositoryTx interface {
	AuthorsCreator
	AuthorsRead
}
