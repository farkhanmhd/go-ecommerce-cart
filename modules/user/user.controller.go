package user

import (
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	"farkhanmhd/go-ecommerce-cart/common/token"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewUserController(
	router *httprouter.Router,
	pool *pgxpool.Pool,
	validate *validator.Validate,
	signer *token.TokenSigner,
) {
	authMiddleware := middleware.NewAuthMiddleware(signer)
	userRepository := NewUserRepository()
	userService := NewUserService(userRepository, pool, validate, signer)
	userHandler := NewUserHandler(userService)

	router.GET("/api/users/current", authMiddleware.Protect(userHandler.FindById))
}
