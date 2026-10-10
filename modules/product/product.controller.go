package product

import (
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	"farkhanmhd/go-ecommerce-cart/common/token"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
)

func NewProductController(
	router *httprouter.Router,
	pool *pgxpool.Pool,
	validate *validator.Validate,
	signer *token.TokenSigner,
) {
	authMiddleware := middleware.NewAuthMiddleware(signer)
	ProductRepository := NewProductRepository()
	ProductService := NewProductService(ProductRepository, pool, validate)
	ProductHandler := NewProductHandler(ProductService)

	router.GET("/api/products", ProductHandler.FindAll)
	router.POST("/api/products", authMiddleware.Protect(ProductHandler.Create))
	router.GET("/api/products/:productId", ProductHandler.FindById)
	router.PUT("/api/products/:productId", authMiddleware.Protect(ProductHandler.Update))
	router.DELETE("/api/products/:productId", authMiddleware.Protect(ProductHandler.Delete))
}
