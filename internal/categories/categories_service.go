package categories

import "context"

type CategoryService interface {
	ICategoryRepository
}

type serviceCategory struct {
	cat ICategoryRepository
}

func NewCategoryService(cat ICategoryRepository) *serviceCategory {
	return &serviceCategory{cat: cat}
}

func (s *serviceCategory) Create(ctx context.Context, c *Category) error {
	if err := c.Validate(); err != nil {
		return err
	}

	return s.cat.Create(ctx, c)
}

func (s *serviceCategory) GetAll(ctx context.Context) ([]Category, error) {
	return s.cat.GetAll(ctx)
}

func (s *serviceCategory) GetById(ctx context.Context, id int64) (*Category, error) {
	category, err := s.cat.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *serviceCategory) Update(ctx context.Context, c *Category) error {
	if err := c.Validate(); err != nil {
		return err
	}

	return s.cat.Update(ctx, c)
}

func (s *serviceCategory) Delete(ctx context.Context, id int64) error {
	err := s.cat.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
