package books

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/authors"
	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/categories"
)

type Filters struct {
	Title    string
	Authors  string
	Category string
	Page     int
}

func (r *BookRepository) GetAll(ctx context.Context, filter *Filters) ([]Books, error) {
	booksMap := make(map[int64]*Books)

	sql := `SELECT b.id, b.title, b.author_id, b.description, b.content,
		b.created_at,b.updated_at,
		c.id, c.name,c.created_at,
		a.id, a.name, a.description
		FROM book_category bc
		JOIN books b ON bc.book_id = b.id
		JOIN categories c ON bc.category_id = c.id
		JOIN authors a ON b.author_id = a.id`

	conditions := []string{}
	params := []interface{}{}

	if filter.Title != "" {
		conditions = append(conditions, "b.title like ?")
		params = append(params, fmt.Sprintf("%s%%", filter.Title))
	}

	if filter.Authors != "" {
		conditions = append(conditions, "a.name like ?")
		params = append(params, fmt.Sprintf("%s%%", filter.Authors))
	}

	if filter.Category != "" {
		conditions = append(conditions, "c.name like ?")
		params = append(params, fmt.Sprintf("%s%%", filter.Category))
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	hasFilter := filter.Title != "" ||
		filter.Authors != "" ||
		filter.Category != ""

	if !hasFilter && filter.Page > 0 {
		size := 10
		offset := (filter.Page - 1) * size
		sql += " LIMIT ? OFFSET ?"
		params = append(params, size, offset)
	}

	rows, err := r.db.QueryContext(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			bookID, categoryId, authorId, IDAuthor      int64
			title, description, content                 string
			categoryName, authorName, authorDec         string
			createdAtRaw, updatedAtRaw, createdAtCatRaw time.Time
		)

		err := rows.Scan(&bookID, &title, &authorId, &description, &content,
			&createdAtRaw, &updatedAtRaw, &categoryId, &categoryName, &createdAtCatRaw,
			&IDAuthor, &authorName, &authorDec,
		)
		if err != nil {
			return nil, err
		}

		createdAt := createdAtRaw.Format("01/01/01 15:04:05")
		updatedAt := updatedAtRaw.Format("01/01/01 15:04:05")
		createdAtCat := createdAtCatRaw.Format("01/01/01 15:04:05")

		if _, ok := booksMap[bookID]; !ok {
			booksMap[bookID] = &Books{
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

		booksMap[bookID].Categories = append(booksMap[bookID].Categories, categories.Category{
			ID:        categoryId,
			Name:      categoryName,
			CreatedAT: createdAtCat,
		})

		booksMap[bookID].Authors = authors.Authors{
			ID:          IDAuthor,
			Name:        authorName,
			Description: authorDec,
		}
	}
	var books []Books
	for _, book := range booksMap {
		books = append(books, *book)
	}

	return books, nil
}
