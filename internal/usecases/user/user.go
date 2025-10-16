package userusecase

import (
	"context"

	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	userrepo "github.com/codepnw/blog-api/internal/repositories/user"
	"github.com/codepnw/blog-api/internal/usecases"
	"github.com/codepnw/blog-api/internal/utils/errs"
	jwttoken "github.com/codepnw/blog-api/internal/utils/jwt"
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

	// Auth
	Register(ctx context.Context, input *userdomain.User) (*AuthResponse, error)
	Login(ctx context.Context, input *userdomain.User) (*AuthResponse, error)
}

type usecase struct {
	repo  userrepo.Repository
	token *jwttoken.JWTToken
}

func NewUserUsecase(repo userrepo.Repository, token *jwttoken.JWTToken) Usecase {
	return &usecase{
		repo:  repo,
		token: token,
	}
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

//  ------- Auth -----------

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (u *usecase) Register(ctx context.Context, input *userdomain.User) (*AuthResponse, error) {
	user, err := u.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	response, err := u.tokenResponse(user)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *usecase) Login(ctx context.Context, input *userdomain.User) (*AuthResponse, error) {
	user, err := u.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errs.ErrUserInvalid
	}

	ok := password.VerifyPassword(input.PasswordHash, user.PasswordHash)
	if !ok {
		return nil, errs.ErrUserInvalid
	}

	response, err := u.tokenResponse(user)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *usecase) tokenResponse(user *userdomain.User) (*AuthResponse, error) {
	accessToken, err := u.token.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.token.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	response := new(AuthResponse)
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	return response, nil
}
