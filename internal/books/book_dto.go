package books

import (
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
)

// @Description Dados necessários para criar um livro
type BookRequest struct {
	Title       string `json:"title" binding:"required" example:"A menina e o porquinho"`
	Description string `json:"description" binding:"required" example:"Esse livro é sobre conteudo infantil"`
	Content     string `json:"content" binding:"required" example:"A menina é o porquinho"`
	AuthorID    int64  `json:"author_id" binding:"required" example:"1"`
}

// @Description Dados necessários pra fazer o relacionamento
type BookCategoryRequest struct {
	BookID     int64 `json:"book_id" binding:"required" example:"1"`
	CategoryID int64 `json:"category_id" binding:"required" example:"1"`
}

type BookResponse struct {
	ID          int64                 `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Content     string                `json:"content"`
	CreatedAt   string                `json:"created_at"`
	UpdatedAt   string                `json:"updated_at"`
	Categories  []categories.Category `json:"categories"`
	AuthorID    int64                 `json:"author_id"`
	Authors     authors.Authors       `json:"author"`
}

func ToResponse(book *Books) BookResponse {
	return BookResponse(*book)
}
