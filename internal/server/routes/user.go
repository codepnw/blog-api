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
	uc := userusecase.NewUserUsecase(repo)
	handler := userhandler.NewUserHandler(uc)

	r := cfg.APP.Group(cfg.Prefix + "/users")
	userID := fmt.Sprintf("/:%s", handlers.ParamKeyUserID)

	r.Post("/", handler.CreateUser)
	r.Get("/", handler.GetAllUsers)
	r.Get(userID, handler.GetUser)
	r.Patch(userID, handler.UpdateUser)
	r.Delete(userID, handler.DeleteUser)
}
