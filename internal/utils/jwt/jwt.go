package jwttoken

import (
	"errors"
	"time"

	"github.com/codepnw/blog-api/internal/config"
	userdomain "github.com/codepnw/blog-api/internal/domains/user"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	secretKey  string
	refreshKey string
}

type UserClaims struct {
	UserID string
	Email  string
	Role   string
	*jwt.RegisteredClaims
}

func InitJWT(cfg *config.EnvConfig) (*JWTToken, error) {
	if cfg == nil {
		return nil, errors.New("jwt config is required")
	}

	if cfg.JWT.SecretKey == "" || cfg.JWT.RefreshKey == "" {
		return nil, errors.New("secret key & refresh key is required")
	}

	return &JWTToken{
		secretKey:  cfg.JWT.SecretKey,
		refreshKey: cfg.JWT.RefreshKey,
	}, nil
}

func (j *JWTToken) GenerateAccessToken(user *userdomain.User) (string, error) {
	duration := time.Hour * 24
	return j.generateToken(j.secretKey, user, duration)
}

func (j *JWTToken) GenerateRefreshToken(user *userdomain.User) (string, error) {
	duration := time.Hour * 24 * 7
	return j.generateToken(j.refreshKey, user, duration)
}

// ---- Generate Token ------

func (j *JWTToken) generateToken(key string, user *userdomain.User, durarion time.Duration) (string, error) {
	claims := &UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(durarion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (j *JWTToken) VerifyAccessToken(token string) (*UserClaims, error) {
	return j.verifyToken(j.secretKey, token)
}

func (j *JWTToken) VerifyRefreshToken(token string) (*UserClaims, error) {
	return j.verifyToken(j.refreshKey, token)
}

// ---- Verify Token ------

func (j *JWTToken) verifyToken(key, tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !token.Valid || !ok {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
