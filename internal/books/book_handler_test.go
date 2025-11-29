package books

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookHandler_CreateBook(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   interface{}
		setupMock func(*MockBookRepo)
		status    int
		body      string
	}{
		{
			name:    "sucesse return 201 Created",
			reqBody: validBookRequest,
			setupMock: func(b *MockBookRepo) {
				b.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
			},
			status: http.StatusCreated,
			body:   expectedBody,
		},
		{
			name:    "BookService Return 500 InternalServerError",
			reqBody: validBookRequest,
			setupMock: func(b *MockBookRepo) {
				b.On("Create", mock.Anything, mock.Anything).Return(errors.New("db timeout")).Once()
			},
			status: http.StatusInternalServerError,
			body:   `{"code":"INTERNAL_ERROR","message":"Erro interno ocorrido","path":"/api/books","status":500}`,
		},
		{
			name:    "InvalidJson Return 400 BadRequest",
			reqBody: `{"Title":"Tx"`,
			setupMock: func(b *MockBookRepo) {
				b.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			},
			status: http.StatusBadRequest,
			body:   `{"code":"BAD_REQUEST","message":"Solicitação invalida","path":"/api/books","status":400}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := new(MockBookRepo)
			router, w := setupTest(mockHandler)

			tt.setupMock(mockHandler)

			reqBoby := createRequest(t, tt.reqBody)

			req, _ := http.NewRequest("POST", "/api/books", reqBoby)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code, "O status deve ser o esperado")

			assert.JSONEq(t, tt.body, w.Body.String())

			mockHandler.AssertExpectations(t)
		})
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	tests := []struct {
		name      string
		bookID    string
		reqBody   interface{}
		setupMock func(*MockBookRepo)
		status    int
		body      string
	}{
		{
			name:    "sucesso Return 204 NoContent",
			bookID:  "1",
			reqBody: updateBookRequest,
			setupMock: func(b *MockBookRepo) {
				b.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			status: http.StatusNoContent,
			body:   "",
		},
		{
			name:    "BookService Return 500 InternalServerError",
			bookID:  "1",
			reqBody: updateBookRequest,
			setupMock: func(b *MockBookRepo) {
				b.On("Update", mock.Anything, mock.Anything).Return(errors.New("db timeout")).Once()
			},
			status: http.StatusInternalServerError,
			body:   `{"code":"INTERNAL_ERROR","message":"Erro interno ocorrido","path":"/api/books/1","status":500}`,
		},
		{
			name:    "InvalidJson Return 400 BadRequest",
			bookID:  "1",
			reqBody: `{"title": "Test`,
			setupMock: func(b *MockBookRepo) {
				b.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
			status: http.StatusBadRequest,
			body:   `{"code":"BAD_REQUEST","message":"Solicitação invalida","path":"/api/books/1","status":400}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := new(MockBookRepo)
			router, w := setupTest(mockHandler)

			tt.setupMock(mockHandler)

			reqBody := createRequest(t, tt.reqBody)

			req, _ := http.NewRequest("PUT", "/api/books/"+tt.bookID, reqBody)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code, "O status deve ser o esperado")

			if tt.body == "" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, tt.body, w.Body.String())
			}
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestBookHandler_DeleteBook(t *testing.T) {
	tests := []struct {
		name      string
		bookID    string
		reqBody   interface{}
		setupMock func(*MockBookRepo)
		status    int
		body      string
	}{
		{
			name:   "secesse Return 204 NoContent",
			bookID: "1",
			setupMock: func(b *MockBookRepo) {
				b.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			status: http.StatusNoContent,
			body:   "",
		},
		{
			name:   "NotFound Return 404",
			bookID: "99",
			setupMock: func(b *MockBookRepo) {
				b.On("Delete", mock.Anything, int64(99)).Return(errors.New("db timeout")).Once()
			},
			status: http.StatusNotFound,
			body:   `{"code":"NOT_FOUND","message":"Recurno não encontrado","path":"/api/books/99","status":404}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBookRepo)
			router, w := setupTest(mockService)

			tt.setupMock(mockService)

			req, _ := http.NewRequest("DELETE", "/api/books/"+tt.bookID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.status, w.Code)
			assert.Contains(t, w.Body.String(), tt.body)
			mockService.AssertExpectations(t)
		})
	}
}
