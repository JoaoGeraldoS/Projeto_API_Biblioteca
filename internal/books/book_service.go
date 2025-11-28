package books

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence/transaction"
)

// type bookCreator interface {
// 	Create(ctx context.Context, b *Books) error
// }

// type bookRead interface {
// 	GetAll(ctx context.Context, filter *Filters)
// }

type BookUsecase interface {
	BookCreator
	BookRead
}

type BookService struct {
	uow  transaction.UnitOfWorkInterface
	auth authors.AuthorRepositoryTx
	cat  categories.CategoryRepositoryTx
	book BookRepositoryTx
}

func NewBookService(u transaction.UnitOfWorkInterface, a authors.AuthorRepositoryTx,
	c categories.CategoryRepositoryTx, b BookRepositoryTx,
) *BookService {
	return &BookService{uow: u, auth: a, cat: c, book: b}
}

func (s *BookService) Create(ctx context.Context, b *Books) error {
	return s.uow.Execute(ctx, func(exec persistence.Executer) error {

		if err := b.Validate(); err != nil {
			return err
		}

		authRepo := s.auth.WithTx(exec)
		catRepo := s.cat.WithTx(exec)
		bookRepo := s.book.WithTx(exec)

		if err := authRepo.Create(ctx, &b.Authors); err != nil {
			return err
		}

		if err := bookRepo.Create(ctx, b); err != nil {
			return err
		}

		for i := range b.Categories {
			c := &b.Categories[i]

			if err := catRepo.Create(ctx, c); err != nil {
				return err
			}

			if err := bookRepo.RelationBookCategory(ctx, b.ID, c.ID); err != nil {
				return err
			}
		}
		return nil
	})
}
