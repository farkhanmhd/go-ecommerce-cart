package auth

import (
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	"farkhanmhd/go-ecommerce-cart/common/token"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewAuthController(
	router *httprouter.Router,
	pool *pgxpool.Pool,
	validate *validator.Validate,
	signer *token.TokenSigner,
) {
	authRepository := NewAuthRepository()
	authService := NewAuthService(authRepository, pool, validate, signer)
	authHandler := NewAuthHandler(authService, signer)
	authMiddleware := middleware.NewAuthMiddleware(signer)

	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)
	router.POST("/api/auth/logout", authMiddleware.Protect(authHandler.Logout))
	router.POST("/api/auth/refresh", authHandler.Refresh)
}
