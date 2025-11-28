package books

import (
	"context"
)

type BookServcie interface {
	BookCreator
	BookRead
}

type serviceBook struct {
	book IBookRepository
}

func NewBookService(b IBookRepository) *serviceBook {
	return &serviceBook{book: b}
}

func (s *serviceBook) Create(ctx context.Context, b *Books) error {

	if err := b.Validate(); err != nil {
		return err
	}

	return s.book.Create(ctx, b)
}

func (s *serviceBook) GetAll(ctx context.Context, filter *Filters) ([]Books, error) {
	return s.book.GetAll(ctx, filter)
}

func (s *serviceBook) GetById(ctx context.Context, id int64) (*Books, error) {
	return s.book.GetById(ctx, id)
}

func (s *serviceBook) RelationBookCategory(ctx context.Context, bookID, catID int64) error {
	return s.book.RelationBookCategory(ctx, bookID, catID)
}
