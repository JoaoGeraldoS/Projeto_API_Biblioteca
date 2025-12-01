package books

import (
	"context"
	"fmt"
	"time"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
)

func (r *BookRepository) GetById(ctx context.Context, id int64) (*Books, error) {
	bookMap := make(map[int64]*Books)
	var book *Books

	sql := `SELECT b.id, b.title, b.author_id, b.description, b.content,
		b.created_at,b.updated_at,
		c.id, c.name,c.created_at,
		a.id, a.name, a.description
		FROM book_category bc
		JOIN books b ON bc.book_id = b.id
		JOIN categories c ON bc.category_id = c.id
		JOIN authors a ON b.author_id = a.id WHERE b.id = ?`

	var (
		bookID, categoryId, authorId, IDAuthor               int64
		title, description, content                          string
		categoryName, authorName, authorDec                  string
		createdAtRawStr, updatedAtRawStr, createdAtCatRawStr string
	)

	err := r.db.QueryRowContext(ctx, sql, id).Scan(&bookID, &title, &authorId, &description, &content,
		&createdAtRawStr, &updatedAtRawStr, &categoryId, &categoryName, &createdAtCatRawStr,
		&IDAuthor, &authorName, &authorDec,
	)
	if err != nil {
		return nil, err
	}

	const timeLayoutDB = "2006-01-02T15:04:05Z"
	createdAtRaw, err := time.Parse(timeLayoutDB, createdAtRawStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar created_at do livro ('%s'): %w", createdAtRawStr, err)
	}

	updatedAtRaw, err := time.Parse(timeLayoutDB, updatedAtRawStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar updated_at do livro ('%s'): %w", updatedAtRawStr, err)
	}

	createdAtCatRaw, err := time.Parse(timeLayoutDB, createdAtCatRawStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao analisar created_at da categoria ('%s'): %w", createdAtCatRawStr, err)
	}

	createdAt := createdAtRaw.Format("01/01/01 15:04:05")
	updatedAt := updatedAtRaw.Format("01/01/01 15:04:05")
	createdAtCat := createdAtCatRaw.Format("01/01/01 15:04:05")

	if _, ok := bookMap[bookID]; !ok {
		bookMap[bookID] = &Books{
			ID:          bookID,
			Title:       title,
			AuthorID:    authorId,
			Description: description,
			Content:     content,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Categories:  []categories.Category{},
			Authors:     authors.Authors{},
		}
	}

	bookMap[bookID].Categories = append(bookMap[bookID].Categories, categories.Category{
		ID:        categoryId,
		Name:      categoryName,
		CreatedAT: createdAtCat,
	})

	bookMap[bookID].Authors = authors.Authors{
		ID:          IDAuthor,
		Name:        authorName,
		Description: authorDec,
	}

	book = bookMap[bookID]

	return book, nil
}
