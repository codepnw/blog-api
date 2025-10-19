package commentusecase

import (
	"context"
	"strings"

	commentdomain "github.com/codepnw/blog-api/internal/domains/comment"
	commentrepo "github.com/codepnw/blog-api/internal/repositories/comment"
	"github.com/codepnw/blog-api/internal/usecases"
	postusecase "github.com/codepnw/blog-api/internal/usecases/post"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
	"github.com/codepnw/blog-api/internal/utils/errs"
	"github.com/codepnw/blog-api/internal/utils/logger"
)

type Usecase interface {
	CreateComment(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error)
	GetCommentByPost(ctx context.Context, postID string) ([]*commentdomain.Comment, error)
	EditComment(ctx context.Context, input *commentdomain.Comment, role string) (*commentdomain.Comment, error)
	DeleteComment(ctx context.Context, commentID int64, userID, role string) error
}

type usecase struct {
	repo commentrepo.Repository
	user userusecase.Usecase
	post postusecase.Usecase
}

func NewCommentUsecase(repo commentrepo.Repository, user userusecase.Usecase, post postusecase.Usecase) Usecase {
	return &usecase{
		repo: repo,
		user: user,
		post: post,
	}
}

func (u *usecase) CreateComment(ctx context.Context, input *commentdomain.Comment) (*commentdomain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	if err := u.validateComment(ctx, input); err != nil {
		return nil, err
	}
	return u.repo.Insert(ctx, input)
}

func (u *usecase) GetCommentByPost(ctx context.Context, postID string) ([]*commentdomain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	return u.repo.ListByPost(ctx, postID)
}

func (u *usecase) EditComment(ctx context.Context, input *commentdomain.Comment, role string) (*commentdomain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	if err := u.validateOwner(ctx, input.ID, input.UserID, role); err != nil {
		logger.Error("usecase.EditComment: validate owner", "id", input.ID, "user_id", input.UserID, "role", role, "error", err)
		return nil, err
	}

	if err := u.validateComment(ctx, input); err != nil {
		return nil, err
	}
	return u.repo.Update(ctx, input)
}

func (u *usecase) DeleteComment(ctx context.Context, commentID int64, userID, role string) error {
	ctx, cancel := context.WithTimeout(ctx, usecases.ContextTimeout)
	defer cancel()

	if err := u.validateOwner(ctx, commentID, userID, role); err != nil {
		logger.Error("usecase.DeleteComment: validate owner", "id", commentID, "user_id", userID, "role", role, "error", err)
		return err
	}
	return u.repo.Delete(ctx, commentID)
}

func (u *usecase) validateComment(ctx context.Context, input *commentdomain.Comment) error {
	_, err := u.post.GetByID(ctx, input.PostID)
	if err != nil {
		logger.Error("usecase.validateComment: get post", "id", input.PostID, "error", err)
		return err
	}

	_, err = u.user.GetUser(ctx, input.UserID)
	if err != nil {
		logger.Error("usecase.validateComment: get user", "id", input.UserID, "error", err)
		return err
	}

	input.Content = strings.TrimSpace(input.Content)
	if input.Content == "" {
		return errs.ErrCommentIsRequired
	}
	return nil
}

func (u *usecase) validateOwner(ctx context.Context, commentID int64, userID, role string) error {
	result, err := u.repo.FindByID(ctx, commentID)
	if err != nil {
		return err
	}

	if role == string(userusecase.RoleAdmin) {
		return nil
	}

	if result.UserID != userID {
		return errs.ErrCommentNotOwner
	}
	return nil
}
