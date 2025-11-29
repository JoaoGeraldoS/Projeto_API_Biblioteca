package books

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookHandler_ReadAllBooks(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockReturn     []Books
		mockErr        error
		expectedStatus int
	}{
		{
			name:        "sucesso_listar_livros",
			queryParams: "?page=1&title=&author=&category=",
			mockReturn: []Books{
				{ID: 1, Title: "Livro 1"},
				{ID: 2, Title: "Livro 2"},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "erro_do_service",
			queryParams:    "?page=1&title=&author=&category=",
			mockReturn:     nil,
			mockErr:        errors.New("erro interno"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "page_invalida",
			queryParams:    "?page=abc",
			mockReturn:     nil,
			mockErr:        nil,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockSvc := new(MockBookRepo)

			if tt.expectedStatus != http.StatusInternalServerError || tt.name != "page_invalida" {
				mockSvc.
					On("GetAll", mock.Anything, mock.Anything).
					Return(tt.mockReturn, tt.mockErr)
			}

			router, rec := setupTest(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/books"+tt.queryParams, nil)
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("status esperado %d, mas veio %d. Body=%s",
					tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestBookHandler_GetBookByID(t *testing.T) {
	tests := []struct {
		name           string
		idParam        string
		mockSetup      func(m *MockBookRepo)
		expectedStatus int
	}{
		{
			name:    "sucesso_buscar_livro",
			idParam: "1",
			mockSetup: func(m *MockBookRepo) {
				m.On("GetById", mock.Anything, int64(1)).
					Return(&Books{ID: 1, Title: "Go"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:    "erro_do_service_retorna_not_found",
			idParam: "2",
			mockSetup: func(m *MockBookRepo) {
				m.On("GetById", mock.Anything, int64(2)).
					Return(nil, errors.New("qualquer erro"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "id_invalido",
			idParam:        "abc",
			mockSetup:      func(m *MockBookRepo) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockBookRepo)
			router, w := setupTest(mockSvc)

			tt.mockSetup(mockSvc)

			req := httptest.NewRequest(http.MethodGet, "/api/books/"+tt.idParam, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBookHandler_RelationBookCategory(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(m *MockBookRepo)
		expectedStatus int
	}{
		{
			name: "sucesso_relacionar",
			body: `{"book_id":1, "category_id":2}`,
			mockSetup: func(m *MockBookRepo) {
				m.On("RelationBookCategory", mock.Anything, int64(1), int64(2)).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "erro_no_bind_json",
			body:           `{"book_id": "abc", "category_id": 2}`,
			mockSetup:      func(m *MockBookRepo) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "erro_do_service",
			body: `{"book_id":1, "category_id":2}`,
			mockSetup: func(m *MockBookRepo) {
				m.On("RelationBookCategory", mock.Anything, int64(1), int64(2)).
					Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockBookRepo)
			tt.mockSetup(mockSvc)

			router, rec := setupTest(mockSvc)

			req, _ := http.NewRequest(http.MethodPost, "/api/books/relation", bytes.NewBuffer([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
