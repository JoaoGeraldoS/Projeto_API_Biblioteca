package books

import (
	"context"
	"time"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
)

func (r *BookRepository) GetById(ctx context.Context, id int64) (*Books, error) {
	bookMap := make(map[int64]*Books)

	sql := `SELECT b.id, b.title, b.author_id, b.description, b.content,
		b.created_at, b.updated_at,
		c.id, c.name, c.created_at,
		a.id, a.name, a.description
		FROM book_category bc
		JOIN books b ON bc.book_id = b.id
		JOIN categories c ON bc.category_id = c.id
		JOIN authors a ON b.author_id = a.id
		WHERE b.id = ?`

	rows, err := r.db.QueryContext(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bookID, categoryId, authorId, authorID      int64
			title, description, content                 string
			categoryName, authorName, authorDesc        string
			createdAtStr, updatedAtStr, createdAtCatStr string
		)

		err := rows.Scan(
			&bookID, &title, &authorId, &description,
			&content, &createdAtStr, &updatedAtStr,
			&categoryId, &categoryName, &createdAtCatStr,
			&authorID, &authorName, &authorDesc,
		)
		if err != nil {
			return nil, err
		}

		createdAt, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			createdAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				return nil, err
			}
		}
		updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
		if err != nil {
			updatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAtStr)
			if err != nil {
				return nil, err
			}
		}
		createdAtCat, err := time.Parse(time.RFC3339, createdAtCatStr)
		if err != nil {
			createdAtCat, err = time.Parse("2006-01-02 15:04:05", createdAtCatStr)
			if err != nil {
				return nil, err
			}
		}

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

		bookMap[bookID].Categories = append(
			bookMap[bookID].Categories,
			categories.Category{
				ID:        categoryId,
				Name:      categoryName,
				CreatedAT: createdAtCat,
			},
		)

		bookMap[bookID].Authors = authors.Authors{
			ID:          authorID,
			Name:        authorName,
			Description: authorDesc,
		}
	}

	for _, book := range bookMap {
		return book, nil
	}

	return nil, nil
}
