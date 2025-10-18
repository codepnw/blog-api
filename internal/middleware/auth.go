package middleware

import (
	"errors"
	"slices"
	"strings"

	"github.com/codepnw/blog-api/internal/handlers"
	jwttoken "github.com/codepnw/blog-api/internal/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

const UserContextKey = "user-context"

type AppMiddleware struct {
	token *jwttoken.JWTToken
}

func InitMiddleware(token *jwttoken.JWTToken) (*AppMiddleware, error) {
	if token == nil {
		return nil, errors.New("token is required")
	}
	return &AppMiddleware{token: token}, nil
}

func (m *AppMiddleware) Authorized() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return errors.New("header is missing")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.New("invalid token format")
		}

		claims, err := m.token.VerifyAccessToken(parts[1])
		if err != nil {
			return err
		}

		ctx.Locals(UserContextKey, claims)
		return ctx.Next()
	}
}

func (m *AppMiddleware) RoleRequired(roles ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userCtx := ctx.Locals(UserContextKey)
		if userCtx == nil {
			return handlers.Unauthorized(ctx, "user context is missing")
		}

		user, ok := userCtx.(*jwttoken.UserClaims)
		if !ok {
			return handlers.Unauthorized(ctx, "invalid user context")
		}

		// for _, role := range roles {
		// 	if user.Role == role {
		// 		return ctx.Next()
		// 	}
		// }
		if slices.Contains(roles, user.Role) {
			return ctx.Next()
		}

		return handlers.Forbidden(ctx, "no permissions")
	}
}

func GetCurrentUser(ctx *fiber.Ctx) (*jwttoken.UserClaims, error) {
	userCtx := ctx.Locals(UserContextKey)
	if userCtx == nil {
		return nil, errors.New("user context is missing")
	}

	user, ok := userCtx.(*jwttoken.UserClaims)
	if !ok {
		return nil, errors.New("invalid user context")
	}
	return user, nil
}
