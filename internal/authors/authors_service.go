package authors

import (
	"context"
)

type AuthorsService interface {
	IAuthorRepository
}

type serviceAuthors struct {
	repo IAuthorRepository
}

func NewAuthorsService(repo IAuthorRepository) *serviceAuthors {
	return &serviceAuthors{repo: repo}
}

func (u *serviceAuthors) Create(ctx context.Context, a *Authors) error {

	if err := a.Validate(); err != nil {
		return err
	}

	return u.repo.Create(ctx, a)
}

func (u *serviceAuthors) GetAll(ctx context.Context) ([]Authors, error) {
	return u.repo.GetAll(ctx)
}

func (u *serviceAuthors) GetByID(ctx context.Context, id int64) (*Authors, error) {
	author, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (u *serviceAuthors) Update(ctx context.Context, a *Authors) error {
	if err := a.Validate(); err != nil {
		return err
	}

	err := u.repo.Update(ctx, a)
	if err != nil {
		return err
	}

	return nil
}

func (u *serviceAuthors) Delete(ctx context.Context, id int64) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
