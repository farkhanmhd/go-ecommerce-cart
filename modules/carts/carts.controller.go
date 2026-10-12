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
	CartHandler := NewCartHandler(CartService)

	router.GET("/api/cart", authMiddleware.Protect(CartHandler.GetCart))
	router.DELETE("/api/cart", authMiddleware.Protect(CartHandler.ClearCart))
	router.POST("/api/cart/items", authMiddleware.Protect(CartHandler.AddItem))
	router.PUT("/api/cart/items/:itemId", authMiddleware.Protect(CartHandler.UpdateItem))
	router.DELETE("/api/cart/items/:itemId", authMiddleware.Protect(CartHandler.RemoveItem))
}
