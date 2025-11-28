package books

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Mock Books
type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) Create(ctx context.Context, b *Books) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBookRepo) RelationBookCategory(ctx context.Context, bookID, catID int64) error {
	args := m.Called(ctx, bookID, catID)
	return args.Error(0)
}

func (m *MockBookRepo) GetAll(ctx context.Context, filter *Filters) ([]Books, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]Books), args.Error(1)
}

func (m *MockBookRepo) GetById(ctx context.Context, id int64) (*Books, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Books), args.Error(1)
}
