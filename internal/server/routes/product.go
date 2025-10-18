package routes

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/handlers"
	posthandler "github.com/codepnw/blog-api/internal/handlers/post"
	postrepo "github.com/codepnw/blog-api/internal/repositories/post"
	postusecase "github.com/codepnw/blog-api/internal/usecases/post"
)

func (cfg *RouteConfig) PostRoutes() {
	repo := postrepo.NewPostRepository(cfg.DB)
	uc := postusecase.NewPostUsecase(repo)
	handler := posthandler.NewPostHandler(uc)

	var (
		basePath     = fmt.Sprintf("%s/posts", cfg.Prefix)
		postIDPath   = fmt.Sprintf("/:%s", handlers.ParamKeyPostID)
		userPostPath = fmt.Sprintf("%s/users/:%s/posts", cfg.Prefix, handlers.ParamKeyAuthorID)
	)

	// Public
	public := cfg.APP.Group(basePath)
	public.Get("/", handler.GetAll)
	public.Get(postIDPath, handler.GetByID)
	// Get By UserID Path
	cfg.APP.Get(userPostPath, handler.GetByUserID)

	// Authorized
	auth := cfg.APP.Group(cfg.Prefix+"/posts", cfg.Mid.Authorized())
	auth.Post("/", handler.Create)
	auth.Patch(postIDPath, handler.Update)
	auth.Delete(postIDPath, handler.Delete)
}
