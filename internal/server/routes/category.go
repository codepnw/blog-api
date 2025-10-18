package routes

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/handlers"
	categoryhandler "github.com/codepnw/blog-api/internal/handlers/category"
	categoryrepo "github.com/codepnw/blog-api/internal/repositories/category"
	categoryusecase "github.com/codepnw/blog-api/internal/usecases/category"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
)

func (cfg *RouteConfig) CategoryRoutes() {
	repo := categoryrepo.NewCategoryRepository(cfg.DB)
	uc := categoryusecase.NewCategoryUsecase(repo)
	handler := categoryhandler.NewCategoryHandler(uc)

	var (
		basePath       = fmt.Sprintf("%s/categories", cfg.Prefix)
		categoryIDPath = fmt.Sprintf("/:%s", handlers.ParamKeyCategoryID)
	)

	// Public
	public := cfg.APP.Group(basePath)
	public.Get("/", handler.GetAll)
	public.Get(categoryIDPath, handler.GetByID)

	// Admin Only
	admin := cfg.APP.Group(basePath, cfg.Mid.Authorized(), cfg.Mid.RoleRequired(string(userusecase.RoleAdmin)))
	admin.Post("/", handler.Create)
	admin.Patch(categoryIDPath, handler.Update)
	admin.Delete(categoryIDPath, handler.Delete)
}
