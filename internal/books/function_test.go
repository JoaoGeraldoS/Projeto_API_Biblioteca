package books

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupTest(mockBookSvc BookServcie) (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	handler := NewBookHandler(mockBookSvc)
	router := gin.New()

	router.Use(middleware.ErrorHandler())

	router.POST("/api/books", handler.CreateBook)
	router.PUT("/api/books/:id", handler.UpdateBook)
	router.DELETE("/api/books/:id", handler.DeleteBook)
	router.GET("/api/books", handler.ReadAllBooks)
	router.GET("/api/books/:id", handler.ReadBook)
	router.POST("/api/books/relation", handler.RelationBookCategory)

	return router, httptest.NewRecorder()
}

var validBookRequest = BookRequest{
	Title:       "A menina e o proquinho",
	Description: "Livro infantil",
	Content:     "A menina e o porquinho",
	AuthorID:    1,
}

var updateBookRequest = BookRequest{
	Title:       "Update Test",
	Description: "Update Test",
	Content:     "Update Test",
	AuthorID:    1,
}

var expectedBody = `{
	"ID":0,
	"Title":"A menina e o proquinho",
	"Description":"Livro infantil",
	"Content":"A menina e o porquinho",
	"CreatedAt":"",
	"UpdatedAt":"",
	"Categories":null,
	"AuthorID":1,
	"Authors":{"ID":0,"Name":"","Description":""}}
`

func createRequest(t *testing.T, data interface{}) io.Reader {
	if jsonString, ok := data.(string); ok {
		return bytes.NewBufferString(jsonString)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Erro ao serializar dados de requisição: %v", err)
	}

	return bytes.NewBuffer(jsonBytes)
}
