package categories

// @Description Dados para criar categoria
type CategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Infantil"`
}

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAT string `json:"created_at"`
}

func ToResponse(cat *Category) CategoryResponse {
	return CategoryResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		CreatedAT: cat.CreatedAT.Format("02/01/06 15:04:05"),
	}
}
