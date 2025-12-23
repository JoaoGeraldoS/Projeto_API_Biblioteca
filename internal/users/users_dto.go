package users

// @Description Dados necessários para criar usuario
type UserRequest struct {
	Name     string `json:"name" binding:"required" example:"Joaquim Silva"`
	Email    string `json:"email" binding:"required" example:"joaquim@email.com"`
	Password string `json:"password" binding:"required" example:"password123"`
	Username string `json:"username" binding:"required" example:"joaoquim324"`
	Role     Roles  `json:"role"`
}

// @Description Dados necessários para atualizar usuario
type UserUpdateRequest struct {
	Name string `json:"name" binding:"required" example:"Joaquim Silva"`
	Bio  string `json:"bio" example:"Meu nome é Joaquim"`
}

// @Description Dados necessários para fazer o login
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"joaquim@email.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Bio       string `json:"bio"`
	Username  string `json:"username"`
	Role      Roles  `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToResponse(u *Users) UserResponse {
	return UserResponse(*u)
}
