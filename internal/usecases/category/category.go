package categoryusecase

import (
	"context"

	categorydomain "github.com/codepnw/blog-api/internal/domains/category"
	categoryrepo "github.com/codepnw/blog-api/internal/repositories/category"
	"github.com/codepnw/blog-api/internal/usecases"
)

type Usecase interface {
	Create(ctx context.Context, input *categorydomain.Category) error
	GetByID(ctx context.Context, id string) (*categorydomain.Category, error)
	GetAll(ctx context.Context) ([]*categorydomain.Category, error)
	Update(ctx context.Context, input *categorydomain.Category) error
	Delete(ctx context.Context, id string) error
}

type usecase struct {
	repo categoryrepo.Repository
}

func NewCategoryUsecase(repo categoryrepo.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(ctx context.Context, input *categorydomain.Category) error {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.Insert(ctx, input)
}

func (u *usecase) GetByID(ctx context.Context, id string) (*categorydomain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.FindByID(ctx, id)
}

func (u *usecase) GetAll(ctx context.Context) ([]*categorydomain.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.List(ctx)
}

func (u *usecase) Update(ctx context.Context, input *categorydomain.Category) error {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.Update(ctx, input)
}

func (u *usecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
