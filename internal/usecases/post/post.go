package postusecase

import (
	"context"
	"time"

	postdomain "github.com/codepnw/blog-api/internal/domains/post"
	postrepo "github.com/codepnw/blog-api/internal/repositories/post"
	"github.com/codepnw/blog-api/internal/utils/logger"
)

const contextTimeout = time.Second * 5

type Usecase interface {
	Create(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error)
	GetByID(ctx context.Context, id string) (*postdomain.Post, error)
	GetByAuthorID(ctx context.Context, authorID string) ([]*postdomain.Post, error)
	GetAll(ctx context.Context) ([]*postdomain.Post, error)
	Update(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error)
	Delete(ctx context.Context, id string) error
}

type usecase struct {
	repo postrepo.Repository
}

func NewPostUsecase(repo postrepo.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) Create(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return u.repo.Insert(ctx, input)
}

func (u *usecase) GetByID(ctx context.Context, id string) (*postdomain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return u.repo.FindByID(ctx, id)
}

func (u *usecase) GetByAuthorID(ctx context.Context, authorID string) ([]*postdomain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return u.repo.FindByAuthorID(ctx, authorID)
}

func (u *usecase) GetAll(ctx context.Context) ([]*postdomain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return u.repo.List(ctx)
}

func (u *usecase) Update(ctx context.Context, input *postdomain.Post) (*postdomain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	_, err := u.repo.FindByID(ctx, input.ID)
	if err != nil {
		logger.Error("usecase.UpdatePost: find post", "id", input.ID, "error", err)
		return nil, err
	}
	return u.repo.Update(ctx, input)
}

func (u *usecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
