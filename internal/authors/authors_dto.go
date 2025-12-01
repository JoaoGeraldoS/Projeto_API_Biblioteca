package authors

// @Description Dados para adicionar um autor
type AuthorRequest struct {
	Name        string `json:"name" binding:"required" example:"João Pereira"`
	Description string `json:"description" binding:"required" example:"Escritor"`
}

type AuthorResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ToResponse(a *Authors) AuthorResponse {
	return AuthorResponse(*a)
}
