package app

import (
	"net/http"

	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/token"
	"farkhanmhd/go-ecommerce-cart/modules/auth"
	"farkhanmhd/go-ecommerce-cart/modules/user"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(pool *pgxpool.Pool,
	validate *validator.Validate,
	signer *token.TokenSigner,
) *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc("GET", "/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.html")
	})

	router.HandlerFunc("GET", "/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "openapi.json")
	})

	auth.NewAuthController(router, pool, validate, signer)
	user.NewUserController(router, pool, validate, signer)

	router.PanicHandler = exceptions.ErrorHandler

	return router
}
