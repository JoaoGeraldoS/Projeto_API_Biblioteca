package books

import (
	"context"
	"sort"
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
		b.created_at, b.updated_at,
		c.id, c.name, c.created_at,
		a.id, a.name, a.description
	FROM book_category bc
	JOIN books b ON bc.book_id = b.id
	JOIN categories c ON bc.category_id = c.id
	JOIN authors a ON b.author_id = a.id`

	var conditions []string
	var params []interface{}

	if filter.Title != "" {
		conditions = append(conditions, "b.title LIKE ?")
		params = append(params, filter.Title+"%")
	}
	if filter.Authors != "" {
		conditions = append(conditions, "a.name LIKE ?")
		params = append(params, filter.Authors+"%")
	}
	if filter.Category != "" {
		conditions = append(conditions, "c.name LIKE ?")
		params = append(params, filter.Category+"%")
	}

	if len(conditions) > 0 {
		sql += " WHERE " + strings.Join(conditions, " AND ")
	}

	sql += " ORDER BY b.id ASC"

	if filter.Page >= 1 {
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
			createdAtCatStr, createdAtStr, updatedAtStr string
		)

		if err := rows.Scan(
			&bookID, &title, &authorId, &description, &content,
			&createdAtStr, &updatedAtStr,
			&categoryId, &categoryName, &createdAtCatStr,
			&IDAuthor, &authorName, &authorDec,
		); err != nil {
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

		booksMap[bookID].Categories = append(
			booksMap[bookID].Categories,
			categories.Category{
				ID:        categoryId,
				Name:      categoryName,
				CreatedAT: createdAtCat,
			},
		)

		booksMap[bookID].Authors = authors.Authors{
			ID:          IDAuthor,
			Name:        authorName,
			Description: authorDec,
		}
	}

	var books []Books
	for _, b := range booksMap {
		books = append(books, *b)
	}

	sort.Slice(books, func(i, j int) bool {
		return books[i].ID < books[j].ID
	})

	return books, nil
}
