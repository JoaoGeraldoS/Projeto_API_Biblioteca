package users

import (
	"context"

	"github.com/JoaoGeraldoS/Projeto_API_Biblioteca/internal/middleware"
)

type userCreator interface {
	Create(ctx context.Context, user *Users) error
	Update(ctx context.Context, user *Users) error
	Delete(ctx context.Context, id int64) error
	Login(ctx context.Context, email, password string) (*Users, error)
}

type userRead interface {
	GetAll(ctx context.Context) ([]Users, error)
	GetById(ctx context.Context, id int64) (*Users, error)
}

type UserService interface {
	userCreator
	userRead
}

type serviceUser struct {
	repo IUsersRepository
}

func NewUsersService(repo IUsersRepository) *serviceUser {
	return &serviceUser{
		repo: repo,
	}
}

func (s *serviceUser) Create(ctx context.Context, user *Users) error {
	if err := user.Validate(); err != nil {
		return err
	}

	password, err := middleware.HashPassowrd(user.Password)
	if err != nil {
		return err
	}

	user.Password = password

	return s.repo.Create(ctx, user)
}

func (s *serviceUser) GetAll(ctx context.Context) ([]Users, error) {
	return s.repo.GetAll(ctx)
}

func (s *serviceUser) GetById(ctx context.Context, id int64) (*Users, error) {
	return s.repo.GetById(ctx, id)
}

func (s *serviceUser) Update(ctx context.Context, user *Users) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *serviceUser) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *serviceUser) Login(ctx context.Context, email, password string) (*Users, error) {
	userDetails, err := s.repo.GetUserDetails(ctx, email)
	if err != nil {
		return nil, err
	}

	if !middleware.VerifyPassword(userDetails.Password, password) {
		return nil, err
	}

	return userDetails, nil
}
