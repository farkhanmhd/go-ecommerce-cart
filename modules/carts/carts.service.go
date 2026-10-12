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
	GetCart(ctx context.Context, userId string) []CartItem
	AddItem(ctx context.Context, userId string, cartItem CartItem) CartItem
	UpdateItem(ctx context.Context, userId string, cartItem CartItem) CartItem
	RemoveItem(ctx context.Context, userId string, itemId string)
	ClearCart(ctx context.Context, userId string)
}

type CartServiceImpl struct {
	CartRepository CartRepository
	Pool           *pgxpool.Pool
	Validate       *validator.Validate
	ProductService product.ProductService
}

func (c *CartServiceImpl) ClearCart(ctx context.Context, userId string) {
	tx, err := c.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	currentCart, err := c.CartRepository.FindByUserId(ctx, tx, userId)

	if err != nil {
		panic(exceptions.NewNotFoundError("CART NOT FOUND"))
	}

	c.CartRepository.DeleteAllItems(ctx, tx, currentCart.Id)
}

func (c *CartServiceImpl) RemoveItem(ctx context.Context, userId string, itemId string) {
	tx, err := c.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	// Check if item is in the cart
	existingItem, err := c.CartRepository.FindItemById(ctx, tx, itemId)

	if err != nil {
		panic(exceptions.NewNotFoundError("EXISTING ITEM NOT FOUND"))
	}

	// check if its owned by the current user
	currentCart, err := c.CartRepository.FindByUserId(ctx, tx, userId)

	if err != nil {
		panic(exceptions.NewNotFoundError("CART NOT FOUND"))
	}

	if currentCart.Id != existingItem.CartId {
		panic(exceptions.NewBadRequestError("BAD REQUEST"))
	}

	c.CartRepository.DeleteItem(ctx, tx, itemId)
}

func (c *CartServiceImpl) UpdateItem(ctx context.Context, userId string, cartItem CartItem) CartItem {
	err := c.Validate.Struct(cartItem)
	utils.PanicIfError(err)
	tx, err := c.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	// Check if item is in the cart
	existingItem, err := c.CartRepository.FindItemById(ctx, tx, cartItem.Id)

	if err != nil {
		panic(exceptions.NewNotFoundError("EXISTING ITEM NOT FOUND"))
	}

	// check if its owned by the current user
	currentCart, err := c.CartRepository.FindByUserId(ctx, tx, userId)

	if err != nil {
		panic(exceptions.NewNotFoundError("CART NOT FOUND"))
	}

	if currentCart.Id != existingItem.CartId {
		panic(exceptions.NewBadRequestError("BAD REQUEST"))
	}

	// Check stock
	product := c.ProductService.FindById(ctx, existingItem.ProductId)

	if cartItem.Quantity > product.Stock {
		panic(exceptions.NewBadRequestError("Quantity cannot be more than stock"))
	}

	cartItem.CartId = currentCart.Id
	return c.CartRepository.UpdateItem(ctx, tx, cartItem)
}

func (c *CartServiceImpl) GetCart(ctx context.Context, userId string) []CartItem {
	tx, err := c.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	// check if user already has a cart. create if not
	cart, err := c.CartRepository.FindByUserId(ctx, tx, userId)
	if err != nil {
		userCart := Cart{UserId: userId}
		cart = c.CartRepository.Save(ctx, tx, userCart)
	}

	return c.CartRepository.FindAllItems(ctx, tx, cart.Id)
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
