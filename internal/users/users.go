package users

import (
	"context"
	"errors"
	"strings"
)

type Roles string

const (
	User  Roles = "user"
	Admin Roles = "admin"
)

type Users struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Bio       string
	Username  string
	Role      Roles
	CreatedAt string
	UpdatedAt string
}

type UserCreator interface {
	Create(ctx context.Context, user *Users) error
	Update(ctx context.Context, user *Users) error
	Delete(ctx context.Context, id int64) error
}

type UserRead interface {
	GetAll(ctx context.Context) ([]Users, error)
	GetById(ctx context.Context, id int64) (*Users, error)
	GetUserDetails(ctx context.Context, email string) (*Users, error)
}

type IUsersRepository interface {
	UserCreator
	UserRead
}

func (u *Users) Validate() error {
	if strings.TrimSpace(u.Name) == "" || strings.TrimSpace(u.Username) == "" {
		return errors.New("nome ou username não podem estar em branco")
	}

	if len(u.Password) < 6 {
		return errors.New("senha deve ser maior que 6 caracteres")
	}

	if u.Role == "" {
		return errors.New("role invaido")
	}

	return nil
}
