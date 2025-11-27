package books_test

import (
	"context"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/database"
)

func TestBookRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *books.Books
		author  *authors.Authors
		wantErr bool
	}{
		{
			name: "sucesso",

			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},

			input: &books.Books{
				Title:       "A menina e o porquinho",
				Description: "Livro a menina e o porquinho",
				Content:     "A menina e o porquinho",
				AuthorID:    1,
			},
			wantErr: false,
		},
		{
			name: "erro dados em branco",
			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},
			input: &books.Books{
				Title:       "",
				Description: "",
				Content:     "",
				AuthorID:    1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := database.SetupTestDB()
			r := books.NewBookRepository(db)
			authRepo := authors.NewAuthorsRepository(db)
			ctx := context.Background()

			authErr := authRepo.Create(ctx, tt.author)

			if authErr != nil {
				t.Fatalf("erro ao criar author")
			}

			err := r.Create(ctx, tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mais não ocorreu: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("não esperava erro, mais ocorreu: %v", err)
			}

			if tt.input.ID == 0 {
				t.Fatalf("ID não foi gerado pela criação")
			}

			var found books.Books
			db.QueryRowContext(ctx, "SELECT id, title FROM books WHERE id = ?", tt.input.ID).
				Scan(&found.ID, &found.Title)

			if found.Title != tt.input.Title {
				t.Errorf("esperava name=%s, recebeu=%s", tt.input.Title, found.Title)
			}

		})
	}
}
