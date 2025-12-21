package categories

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCategoryRepo struct {
	mock.Mock
}

func (m *MockCategoryRepo) Create(ctx context.Context, cat *Category) error {
	args := m.Called(ctx, cat)
	return args.Error(0)
}

func (m *MockCategoryRepo) GetAll(ctx context.Context) ([]Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Category), args.Error(1)
}

func (m *MockCategoryRepo) GetById(ctx context.Context, id int64) (*Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	cat, ok := args.Get(0).(*Category)
	if !ok {
		return nil, args.Error(1)
	}

	return cat, args.Error(1)
}

func (m *MockCategoryRepo) Update(ctx context.Context, b *Category) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
