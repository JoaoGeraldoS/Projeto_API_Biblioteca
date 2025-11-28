package books

import (
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
)

type BookRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	AuthorID    int64  `json:"author_id"`
}

type BookCategoryRequest struct {
	BookID     int64 `json:"book_id"`
	CategoryID int64 `json:"category_id"`
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
