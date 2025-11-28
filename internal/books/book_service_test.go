package books_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookService_Create(t *testing.T) {
	tests := []struct {
		name        string
		input       *books.Books
		authErr     error
		catErr      error
		relationErr error
		bookErr     error
		wantErr     bool
	}{
		{
			name: "sucesso",
			input: &books.Books{
				ID:          1,
				Title:       "A menina e o porquinho",
				Description: "Livro infantil",
				Content:     "A menina e o porquinho",
				AuthorID:    1,
				Authors: authors.Authors{
					ID:          1,
					Name:        "Nao sei",
					Description: "Nao sei",
				},
				Categories: []categories.Category{
					{ID: 1, Name: "Infantil"},
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao criar author",
			input: &books.Books{
				Title:       "Valid Title",
				Description: "Valid Description",
				Authors: authors.Authors{
					Name: "Valid Name",
				},
			},
			authErr: errors.New("erro no author"),
			wantErr: true,
		},
		{
			name: "erro ao criar book",
			input: &books.Books{
				Title:       "Valid Title",
				Description: "Valid Description",
				Authors: authors.Authors{
					Name: "Valid Name",
				},
			},
			bookErr: errors.New("erro no book"),
			wantErr: true,
		},
		{
			name: "erro ao criar categoria",
			input: &books.Books{
				Title:       "Valid Title",
				Description: "Valid Description",
				Authors: authors.Authors{
					Name: "Valid Name",
				},
				Categories: []categories.Category{{Name: "Tech"}},
			},
			catErr:  errors.New("category error"),
			wantErr: true,
		},
		{
			name: "erro ao criar relacao",
			input: &books.Books{
				ID:          1,
				Title:       "Valid Title",
				Description: "Valid Description",
				Authors: authors.Authors{
					Name: "Valid Name",
				},
				Categories: []categories.Category{{ID: 9, Name: "Valid Cat"}},
			},
			relationErr: errors.New("relation error"),
			wantErr:     true,
		},
		{
			name: "titulo em branco",
			input: &books.Books{
				Title:       "",
				Description: "Valid Description",
				Authors: authors.Authors{
					Name: "Valid Name",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUow := new(books.MockUoW)
			mockAuth := new(books.AuthMockRepo)
			mockCat := new(books.CatMockRepo)
			mockBook := new(books.MockBookRepo)

			mockUow.On("Execute", mock.Anything, mock.Anything).
				Return(func(_ context.Context, fn func(persistence.Executer) error) error {
					return fn(nil)
				})

			if tt.name != "titulo em branco" {
				mockAuth.On("WithTx", mock.Anything).Return(mockAuth)
				mockCat.On("WithTx", mock.Anything).Return(mockCat)
				mockBook.On("WithTx", mock.Anything).Return(mockBook)

				mockAuth.On("Create", mock.Anything, mock.MatchedBy(func(a interface{}) bool {
					_, ok := a.(*authors.Authors)
					return ok
				})).Return(tt.authErr)

				if tt.authErr == nil {
					mockBook.On("Create", mock.Anything, mock.MatchedBy(func(b interface{}) bool {
						_, ok := b.(*books.Books)
						return ok
					})).Return(tt.bookErr)

					if tt.bookErr == nil && len(tt.input.Categories) > 0 {
						for _, c := range tt.input.Categories {
							c := c
							mockCat.On("Create", mock.Anything, mock.MatchedBy(func(x interface{}) bool {
								cc, ok := x.(*categories.Category)
								return ok && (cc.Name == c.Name || cc.ID == c.ID)
							})).Return(tt.catErr)

							if tt.catErr == nil {
								mockBook.On("RelationBookCategory", mock.Anything, tt.input.ID, c.ID).Return(tt.relationErr)
							}
						}
					}
				}
			}

			svc := books.NewBookService(mockUow, mockAuth, mockCat, mockBook)

			err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockUow.AssertExpectations(t)
			mockAuth.AssertExpectations(t)
			mockCat.AssertExpectations(t)
			mockBook.AssertExpectations(t)
		})
	}
}
