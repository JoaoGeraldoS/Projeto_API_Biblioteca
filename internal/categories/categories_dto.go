package categories

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAT string `json:"created_at"`
}

func ToResponse(cat *Category) CategoryResponse {
	return CategoryResponse(*cat)
}
