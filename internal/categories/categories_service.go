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
