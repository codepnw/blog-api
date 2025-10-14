package userusecase

import (
	"context"

	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	userrepo "github.com/codepnw/blog-api/internal/repositories/user"
	"github.com/codepnw/blog-api/internal/usecases"
	"github.com/codepnw/blog-api/internal/utils/errs"
	"github.com/codepnw/blog-api/internal/utils/password"
)

type userRole string

const (
	RoleAdmin userRole = "admin"
	RoleUser  userRole = "user"
)

type Usecase interface {
	CreateUser(ctx context.Context, input *userdomain.User) (*userdomain.User, error)
	GetUser(ctx context.Context, id string) (*userdomain.User, error)
	GetAllUsers(ctx context.Context) ([]*userdomain.User, error)
	UpdateUser(ctx context.Context, input *userdomain.User) (*userdomain.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type usecase struct {
	repo userrepo.Repository
}

func NewUserUsecase(repo userrepo.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) CreateUser(ctx context.Context, input *userdomain.User) (*userdomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	hashed, err := password.HashedPassword(input.PasswordHash)
	if err != nil {
		return nil, err
	}

	input.PasswordHash = hashed
	input.Role = string(RoleUser)

	return u.repo.Insert(ctx, input)
}

func (u *usecase) GetAllUsers(ctx context.Context) ([]*userdomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.List(ctx)
}

func (u *usecase) GetUser(ctx context.Context, id string) (*userdomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.FindByID(ctx, id)
}

func (u *usecase) UpdateUser(ctx context.Context, input *userdomain.User) (*userdomain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	user, err := u.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if user.ID != input.ID {
		return nil, errs.ErrUserUnauthorized
	}

	return u.repo.Update(ctx, input)
}

func (u *usecase) DeleteUser(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.Delete(ctx, id)
}
