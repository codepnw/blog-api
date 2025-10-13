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

	r := cfg.APP.Group(cfg.Prefix + "/posts")
	postID := fmt.Sprintf("/:%s", handlers.ParamKeyPostID)
	authorID := fmt.Sprintf("/author/:%s", handlers.ParamKeyAuthorID)

	r.Post("/", handler.Create)
	r.Get("/", handler.GetAll)
	r.Get(postID, handler.GetByID)
	r.Patch(postID, handler.Update)
	r.Delete(postID, handler.Delete)
	r.Get(authorID, handler.GetByAuthorID)
}
