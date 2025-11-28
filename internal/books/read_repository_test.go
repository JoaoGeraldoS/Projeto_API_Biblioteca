package books_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/database"
)

func seedDataUnic(t *testing.T, db *sql.DB) {
	ctx := context.Background()

	authorRepo := authors.NewAuthorsRepository(db)
	categoryRepo := categories.NewCategoryRepository(db)
	bookRepo := books.NewBookRepository(db)

	cat := &categories.Category{ID: 1, Name: "Programação"}

	if err := authorRepo.Create(ctx, &authors.Authors{ID: 1, Name: "Autor X"}); err != nil {
		t.Fatalf("author: %v", err)
	}

	if err := categoryRepo.Create(ctx, cat); err != nil {
		t.Fatalf("category: %v", err)
	}

	if err := bookRepo.Create(ctx, &books.Books{ID: 3, Title: "Go Lang", Description: "D1", Content: "C1", AuthorID: 1}); err != nil {
		t.Fatalf("book: %v", err)
	}

	if err := bookRepo.RelationBookCategory(ctx, 1, cat.ID); err != nil {
		t.Fatalf("relation: %v", err)
	}

}

func TestBookRepository_GetById(t *testing.T) {
	db := database.SetupTestDB()
	seedDataUnic(t, db)
	repo := books.NewBookRepository(db)

	tests := []struct {
		name          string
		wantCount     int
		wantErr       bool
		wantFirstName string
	}{
		{
			name:      "nenhum resultado",
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:          "sucesso",
			wantCount:     1,
			wantFirstName: "Go Lang",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := repo.GetById(context.Background(), 1)

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if !tt.wantErr {
				if got.Title != tt.wantFirstName {
					t.Fatalf("esperava o livro: %v mais veio %v", tt.wantFirstName, got.Title)
				}
			}

			if err != nil {
				t.Fatalf("teve erro ao buscar no banco: %v", err)
			}

			if len(got.Categories) == 0 {
				t.Errorf("esperava categoria")
			}
			if got.AuthorID == 0 {
				t.Errorf("esperava autor")
			}
		})
	}
}
