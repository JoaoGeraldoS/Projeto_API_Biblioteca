package books_test

import (
	"context"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/books"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/database"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/infra/persistence"
)

func seedData(t *testing.T, exec persistence.Executer) {
	ctx := context.Background()

	authorRepo := authors.NewAuthorsRepository(exec)
	categoryRepo := categories.NewCategoryRepository(exec)
	bookRepo := books.NewBookRepository(exec)

	if err := authorRepo.Create(ctx, &authors.Authors{ID: 1, Name: "Autor X"}); err != nil {
		t.Fatalf("author: %v", err)
	}

	if err := categoryRepo.Create(ctx, &categories.Category{ID: 1, Name: "Programação"}); err != nil {
		t.Fatalf("category: %v", err)
	}

	if err := bookRepo.Create(ctx, &books.Books{ID: 1, Title: "Go Lang", Description: "D1", Content: "C1", AuthorID: 1}); err != nil {
		t.Fatalf("book: %v", err)
	}
	if err := bookRepo.Create(ctx, &books.Books{ID: 2, Title: "Python", Description: "D1", Content: "C1", AuthorID: 1}); err != nil {
		t.Fatalf("book: %v", err)
	}

	if err := bookRepo.RelationBookCategory(ctx, 1, 1); err != nil {
		t.Fatalf("relation: %v", err)
	}
	if err := bookRepo.RelationBookCategory(ctx, 2, 1); err != nil {
		t.Fatalf("relation: %v", err)
	}
}

func TestBookRepository_GetAll_TableDriven(t *testing.T) {
	db := database.SetupTestDB()
	seedData(t, db)
	repo := books.NewBookRepository(db)

	tests := []struct {
		name          string
		filter        *books.Filters
		wantCount     int
		wantFirstID   int64
		wantFirstName string
	}{
		{
			name:          "sem filtros",
			filter:        &books.Filters{},
			wantCount:     2,
			wantFirstID:   1,
			wantFirstName: "Go Lang",
		},
		{
			name: "filtro por título",
			filter: &books.Filters{
				Title: "Go",
			},
			wantCount:     1,
			wantFirstID:   1,
			wantFirstName: "Go Lang",
		},
		{
			name: "filtro por autor",
			filter: &books.Filters{
				Authors: "Autor X",
			},
			wantCount:     2,
			wantFirstID:   1,
			wantFirstName: "Go Lang",
		},
		{
			name: "filtro por categoria",
			filter: &books.Filters{
				Category: "Programação",
			},
			wantCount:     2,
			wantFirstID:   1,
			wantFirstName: "Go Lang",
		},
		{
			name: "nenhum resultado",
			filter: &books.Filters{
				Title: "NãoExiste",
			},
			wantCount: 0,
		},
		{
			name: "página ignorada quando WHERE existe",
			filter: &books.Filters{
				Title: "Go",
				Page:  3,
			},
			wantCount:     1,
			wantFirstID:   1,
			wantFirstName: "Go Lang",
		},
	}

	ctx := context.Background()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			got, err := repo.GetAll(ctx, tc.filter)
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if len(got) != tc.wantCount {
				t.Fatalf("qtde errada: esperava %d, veio %d", tc.wantCount, len(got))
			}

			if tc.wantCount == 0 {
				return
			}

			if got[0].ID != tc.wantFirstID {
				t.Errorf("ID errado. esperado %d, veio %d", tc.wantFirstID, got[0].ID)
			}

			if got[0].Title != tc.wantFirstName {
				t.Errorf("Título errado. esperado %s, veio %s", tc.wantFirstName, got[0].Title)
			}

			if len(got[0].Categories) == 0 {
				t.Errorf("esperava categoria")
			}
			if len(got[0].Authors) == 0 {
				t.Errorf("esperava autor")
			}
		})
	}
}
