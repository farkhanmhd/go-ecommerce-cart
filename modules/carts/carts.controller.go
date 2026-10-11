package carts

import (
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	"farkhanmhd/go-ecommerce-cart/common/token"
	"farkhanmhd/go-ecommerce-cart/modules/product"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewCartController(
	router *httprouter.Router,
	pool *pgxpool.Pool,
	validate *validator.Validate,
	signer *token.TokenSigner,
) {
	authMiddleware := middleware.NewAuthMiddleware(signer)
	CartRepository := NewCartRepository()
	ProductRepository := product.NewProductRepository()
	ProductService := product.NewProductService(ProductRepository, pool, validate)
	CartService := NewCartService(CartRepository, pool, validate, ProductService)
	CartNewCartHandler := NewCartHandler(CartService)

	router.POST("/api/cart/items", authMiddleware.Protect(CartNewCartHandler.AddItem))
}
