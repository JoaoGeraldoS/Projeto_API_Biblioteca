package books_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/database"
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
			if err := db.QueryRowContext(ctx, "SELECT id, title FROM books WHERE id = ?", tt.input.ID).
				Scan(&found.ID, &found.Title); err != nil {
				t.Fatalf("erro ao escanear: %v", err)
			}

			if found.Title != tt.input.Title {
				t.Errorf("esperava name=%s, recebeu=%s", tt.input.Title, found.Title)
			}
		})
	}
}

func TestBookRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		author  *authors.Authors
		input   *books.Books
		update  *books.Books
		wantErr bool
	}{
		{
			name: "sucesso",
			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},
			input: &books.Books{
				Title:       "Título Original",
				Description: "Descrição Original",
				Content:     "Conteúdo Original",
				AuthorID:    1,
			},
			update: &books.Books{
				ID:          0,
				Title:       "Título Atualizado",
				Description: "Descrição Atualizada",
				Content:     "Conteúdo Atualizado",
				AuthorID:    1,
			},
			wantErr: false,
		},
		{
			name: "erro livro inexistente",
			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},
			input: nil,
			update: &books.Books{
				ID:          999,
				Title:       "Título Qualquer",
				Description: "Descrição Qualquer",
				Content:     "Conteúdo Qualquer",
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
				t.Fatalf("erro ao criar author: %v", authErr)
			}

			if tt.input != nil {
				err := r.Create(ctx, tt.input)
				if err != nil {
					t.Fatalf("erro ao criar livro inicial: %v", err)
				}
				tt.update.ID = tt.input.ID
			}

			err := r.Update(ctx, tt.update)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas não ocorreu: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("não esperava erro, mas ocorreu: %v", err)
			}

			var found books.Books
			if err := db.QueryRowContext(ctx, "SELECT id, title FROM books WHERE id = ?", tt.update.ID).
				Scan(&found.ID, &found.Title); err != nil {
				t.Fatalf("erro ao escanear: %v", err)
			}

			if found.Title != tt.update.Title {
				t.Errorf("esperava title=%s, recebeu=%s", tt.update.Title, found.Title)
			}

		})
	}
}

func TestBookRepository_Delete(t *testing.T) {
	tests := []struct {
		name     string
		author   *authors.Authors
		input    *books.Books
		deleteID int64
		wantErr  bool
	}{
		{
			name: "sucesso",
			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},
			input: &books.Books{
				Title:       "Título Original",
				Description: "Descrição Original",
				Content:     "Conteúdo Original",
				AuthorID:    1,
			},
			deleteID: 0,
			wantErr:  false,
		},
		{
			name: "erro livro inexistente",
			author: &authors.Authors{
				ID:   1,
				Name: "Teste",
			},
			input:    nil,
			deleteID: 999,
			wantErr:  true,
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
				t.Fatalf("erro ao criar author: %v", authErr)
			}

			if tt.input != nil {
				err := r.Create(ctx, tt.input)
				if err != nil {
					t.Fatalf("erro ao criar livro inicial: %v", err)
				}
				tt.deleteID = int64(tt.input.ID)
			}

			err := r.Delete(ctx, tt.deleteID)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas não ocorreu: %v", err)
				}
				expectedErr := errors.New("erro ao deletar livro")
				if err.Error() != expectedErr.Error() {
					t.Errorf("erro inesperado: esperado %v, recebeu %v", expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("não esperava erro, mas ocorreu: %v", err)
			}
		})
	}
}
