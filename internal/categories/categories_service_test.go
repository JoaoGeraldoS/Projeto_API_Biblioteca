package categories

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_serviceCategory_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *Category
		catErr  error
		wantErr bool
	}{
		{
			name: "sucesso",
			input: &Category{
				Name: "Teste",
			},
			wantErr: false,
		},
		{
			name: "name em branco",
			input: &Category{
				Name: "",
			},
			catErr:  errors.New("nome em branco"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCategoryRepo)

			if tt.name != "name em branco" {
				mockRepo.On("Create", mock.Anything, tt.input).Return(tt.catErr)
			}

			svc := NewCategoryService(mockRepo)

			err := svc.Create(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
