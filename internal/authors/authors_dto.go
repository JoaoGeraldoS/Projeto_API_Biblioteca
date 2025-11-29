package authors

type AuthorRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AuthorResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ToResponse(a *Authors) AuthorResponse {
	return AuthorResponse(*a)
}
