package routes

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/handlers"
	categoryhandler "github.com/codepnw/blog-api/internal/handlers/category"
	categoryrepo "github.com/codepnw/blog-api/internal/repositories/category"
	categoryusecase "github.com/codepnw/blog-api/internal/usecases/category"
)

func (cfg *RouteConfig) CategoryRoutes() {
	repo := categoryrepo.NewCategoryRepository(cfg.DB)
	uc := categoryusecase.NewCategoryUsecase(repo)
	handler := categoryhandler.NewCategoryHandler(uc)

	r := cfg.APP.Group(cfg.Prefix + "/categories")
	categoryID := fmt.Sprintf("/:%s", handlers.ParamKeyCategoryID)

	r.Post("/", handler.Create)
	r.Get("/", handler.GetAll)
	r.Get(categoryID, handler.GetByID)
	r.Patch(categoryID, handler.Update)
	r.Delete(categoryID, handler.Delete)
}
