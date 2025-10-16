package routes

import (
	"fmt"

	"github.com/codepnw/blog-api/internal/handlers"
	userhandler "github.com/codepnw/blog-api/internal/handlers/user"
	userrepo "github.com/codepnw/blog-api/internal/repositories/user"
	userusecase "github.com/codepnw/blog-api/internal/usecases/user"
)

func (cfg *RouteConfig) UserRoutes() {
	repo := userrepo.NewUserRepository(cfg.DB)
	uc := userusecase.NewUserUsecase(repo, cfg.Token)
	handler := userhandler.NewUserHandler(uc)

	// Public
	public := cfg.APP.Group(cfg.Prefix + "/auth")
	public.Post("/register", handler.Register)
	public.Post("/login", handler.Login)

	// Private
	private := cfg.APP.Group(cfg.Prefix+"/users", cfg.Mid.Authorized())
	private.Get("/me", handler.GetProfile)

	// Admin Only
	admin := cfg.APP.Group(cfg.Prefix+"/users", cfg.Mid.Authorized(), cfg.Mid.RoleRequired(string(userusecase.RoleAdmin)))
	userID := fmt.Sprintf("/:%s", handlers.ParamKeyUserID)

	admin.Post("/", handler.CreateUser)
	admin.Get("/", handler.GetAllUsers)
	admin.Get(userID, handler.GetUser)
	admin.Patch(userID, handler.UpdateUser)
	admin.Delete(userID, handler.DeleteUser)
}
