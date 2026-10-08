package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"farkhanmhd/go-ecommerce-cart/common/app"
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	"farkhanmhd/go-ecommerce-cart/common/token"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	pool := app.NewPGPool(os.Getenv("DATABASE_URL"))
	validate := validator.New()

	key, _ := base64.StdEncoding.DecodeString(os.Getenv("JWT_AUTH_SECRET"))
	signer := token.NewTokenSign(key, os.Getenv("JWT_ISSUER"), os.Getenv("JWT_AUDIENCE"), 15*time.Minute)

	router := app.NewRouter(pool, validate, signer)

	server := http.Server{
		Addr:    "localhost:3003",
		Handler: middleware.NewCorsMiddleware(router),
	}

	fmt.Println("Server is now running at http://localhost:3003")

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
