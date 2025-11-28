package books

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
	"github.com/stretchr/testify/mock"
)

type AuthMockRepo struct {
	mock.Mock
}

// Mock authors
func (m *AuthMockRepo) Create(ctx context.Context, a *authors.Authors) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}

func (m *AuthMockRepo) GetAll(ctx context.Context) ([]authors.Authors, error) {
	args := m.Called(ctx)
	return args.Get(0).([]authors.Authors), args.Error(1)
}

func (m *AuthMockRepo) WithTx(exec persistence.Executer) authors.AuthorRepositoryTx {
	args := m.Called(exec)
	return args.Get(0).(authors.AuthorRepositoryTx)
}

// Mock Categories
type CatMockRepo struct {
	mock.Mock
}

func (m *CatMockRepo) Create(ctx context.Context, c *categories.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CatMockRepo) GetAll(ctx context.Context) ([]categories.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]categories.Category), args.Error(1)
}

func (m *CatMockRepo) WithTx(exec persistence.Executer) categories.CategoryRepositoryTx {
	args := m.Called(exec)
	return args.Get(0).(categories.CategoryRepositoryTx)
}

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

func (m *MockBookRepo) WithTx(exec persistence.Executer) BookRepositoryTx {
	args := m.Called(exec)
	return args.Get(0).(BookRepositoryTx)
}

// Mock Transaction
type MockUoW struct {
	mock.Mock
}

func (m *MockUoW) Execute(ctx context.Context, fn func(exec persistence.Executer) error) error {
	args := m.Called(ctx, fn)

	if fnRet, ok := args.Get(0).(func(context.Context, func(persistence.Executer) error) error); ok {
		return fnRet(ctx, fn)
	}

	return args.Error(0)
}
