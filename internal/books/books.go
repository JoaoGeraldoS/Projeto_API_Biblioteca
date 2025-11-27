package books

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

type Books struct {
	ID          int64
	Title       string
	Description string
	Content     string
	CreatedAt   string
	UpdatedAt   string
	Categories  []*categories.Category
	AuthorID    int64
	Authors     []authors.Authors
}

type Uow interface {
	Executer(ctx context.Context, fn func(exec persistence.Executer) error) error
}

type BookCreator interface {
	Create(ctx context.Context, b *Books) error
	WithTx(exec persistence.Executer) BookRepositoryTx
	RelationBookCategory(ctx context.Context, bookID, categoryID int64) error
}

type BookRead interface {
	GetAll(ctx context.Context, filter *Filters) ([]Books, error)
	GetById(ctx context.Context, id int64) (*Books, error)
}

type BookRepositoryTx interface {
	BookCreator
	BookRead
}
