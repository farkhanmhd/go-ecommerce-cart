package carts

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"farkhanmhd/go-ecommerce-cart/modules/product"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartService interface {
	AddItem(ctx context.Context, userId string, cartItem CartItem) CartItem
}

type CartServiceImpl struct {
	CartRepository CartRepository
	Pool           *pgxpool.Pool
	Validate       *validator.Validate
	ProductService product.ProductService
}

// AddItem implements [CartService].
func (c *CartServiceImpl) AddItem(ctx context.Context, userId string, cartItem CartItem) CartItem {
	err := c.Validate.Struct(cartItem)
	utils.PanicIfError(err)
	tx, err := c.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	// check if product exist and has stock
	product := c.ProductService.FindById(ctx, cartItem.ProductId)

	if cartItem.Quantity > product.Stock {
		panic(exceptions.NewBadRequestError("Quantity cannot be more than stock"))
	}

	// check if user already has a cart. create if not
	cart, err := c.CartRepository.FindByUserId(ctx, tx, userId)
	if err != nil {
		userCart := Cart{UserId: userId}
		cart = c.CartRepository.Save(ctx, tx, userCart)
	}

	cartItem.CartId = cart.Id

	// Check if product is already in the cart
	currentCartItem, err := c.CartRepository.FindItemByCartAndProduct(ctx, tx, cart.Id, cartItem.ProductId)
	if err != nil {
		// Insert new item
		return c.CartRepository.SaveItem(ctx, tx, cartItem)
	}

	// check total quantity against available stock
	if currentCartItem.Quantity+cartItem.Quantity > product.Stock {
		panic(exceptions.NewBadRequestError("Total quantity cannot be more than stock"))
	}

	currentCartItem.Quantity += cartItem.Quantity
	return c.CartRepository.UpdateItem(ctx, tx, currentCartItem)
}

func NewCartService(cartRepository CartRepository, pool *pgxpool.Pool, validate *validator.Validate, productService product.ProductService) CartService {
	return &CartServiceImpl{
		CartRepository: cartRepository,
		Pool:           pool,
		Validate:       validate,
		ProductService: productService,
	}
}
