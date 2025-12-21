package categories_test

import (
	"context"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/database"
)

func TestCategoryRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *categories.Category
		wantErr bool
	}{
		{
			name: "sucesso",
			input: &categories.Category{
				Name: "Infantil",
			},
			wantErr: false,
		},
		{
			name: "erro criacao",
			input: &categories.Category{
				Name: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := database.SetupTestDB()
			r := categories.NewCategoryRepository(db)
			ctx := context.Background()

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

			var found categories.Category
			if err := db.QueryRowContext(ctx, "SELECT id, name FROM categories WHERE id = ?", tt.input.ID).
				Scan(&found.ID, &found.Name); err != nil {
				t.Fatalf("erro ao escanear: %v", err)
			}

			if found.Name != tt.input.Name {
				t.Errorf("esperava name=%s, recebeu=%s", tt.input.Name, found.Name)
			}
		})
	}
}
