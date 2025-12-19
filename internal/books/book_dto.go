package books

import (
	"time"

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
	ID          int64                         `json:"id"`
	Title       string                        `json:"title"`
	Description string                        `json:"description"`
	Content     string                        `json:"content"`
	CreatedAt   string                        `json:"created_at"`
	UpdatedAt   string                        `json:"updated_at"`
	Categories  []categories.CategoryResponse `json:"categories"`
	AuthorID    int64                         `json:"author_id"`
	Authors     authors.AuthorResponse        `json:"author"`
}

func toCategoryResponse(cats []categories.Category) []categories.CategoryResponse {
	resp := make([]categories.CategoryResponse, 0, len(cats))
	for _, c := range cats {
		resp = append(resp, categories.ToResponse(&c))
	}
	return resp
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02/01/06 15:04:05")
}

func ToResponse(b *Books) BookResponse {

	return BookResponse{
		ID:          b.ID,
		Title:       b.Title,
		Description: b.Description,
		Content:     b.Content,
		AuthorID:    b.AuthorID,
		CreatedAt:   formatTime(b.CreatedAt),
		UpdatedAt:   formatTime(b.UpdatedAt),
		Authors:     authors.ToResponse(&b.Authors),
		Categories:  toCategoryResponse(b.Categories),
	}
}
