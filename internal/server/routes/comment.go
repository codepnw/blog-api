package routes

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/handlers"
	commenthandler "github.com/codepnw/blog-api/internal/handlers/comment"
	commentrepo "github.com/codepnw/blog-api/internal/repositories/comment"
	postrepo "github.com/codepnw/blog-api/internal/repositories/post"
	userrepo "github.com/codepnw/blog-api/internal/repositories/user"
	commentusecase "github.com/codepnw/blog-api/internal/usecases/comment"
	postusecase "github.com/codepnw/blog-api/internal/usecases/post"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
)

func (cfg *RouteConfig) CommentRoutes() {
	// User Usecase
	userRepo := userrepo.NewUserRepository(cfg.DB)
	userUc := userusecase.NewUserUsecase(userRepo, cfg.Token)

	// Post Usecase
	postRepo := postrepo.NewPostRepository(cfg.DB)
	postUc := postusecase.NewPostUsecase(postRepo)

	// Comment
	repo := commentrepo.NewCommentRepository(cfg.DB)
	uc := commentusecase.NewCommentUsecase(repo, userUc, postUc)
	handler := commenthandler.NewCommentHandler(uc)

	var (
		basePath      = fmt.Sprintf("%s/posts/:%s/comments", cfg.Prefix, handlers.ParamKeyPostID)
		commentIDPath = fmt.Sprintf("/:%s", handlers.ParamKeyCommentID)
	)

	// Public
	cfg.APP.Get(basePath, handler.GetCommentByPost)

	// Private
	private := cfg.APP.Group(basePath, cfg.Mid.Authorized())
	private.Post("/", handler.CreateComment)
	private.Patch(commentIDPath, handler.EditComment)
	private.Delete(commentIDPath, handler.DeleteComment)
}
