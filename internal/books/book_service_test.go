package books_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
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
			},
			wantErr: false,
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
			mockBook := new(books.MockBookRepo)

			if tt.name != "titulo em branco" {
				mockBook.On("Create", mock.Anything, mock.MatchedBy(func(b interface{}) bool {
					_, ok := b.(*books.Books)
					return ok
				})).Return(tt.bookErr)
			}

			svc := books.NewBookService(mockBook)

			err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockBook.AssertExpectations(t)
		})
	}
}

func TestBookService_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		filter   *books.Filters
		books    []books.Books
		wantCont int
		wantErr  bool
		want     string
	}{
		{
			name: "sucesso",
			filter: &books.Filters{
				Title:    "",
				Authors:  "",
				Category: "",
				Page:     10,
			},
			books:    []books.Books{{ID: 1, Title: "T1"}, {ID: 2, Title: "T2"}},
			wantCont: 2,
			wantErr:  false,
		},
		{
			name:     "sem resultado",
			wantCont: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBook := new(books.MockBookRepo)

			mockBook.On("GetAll", mock.Anything, tt.filter).Return(tt.books, nil)

			svc := books.NewBookService(mockBook)
			result, err := svc.GetAll(context.Background(), tt.filter)

			if tt.wantErr {
				t.Fatalf("erro ao buscar dados")
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.books, result)

			mockBook.AssertExpectations(t)
		})
	}
}
