package carts

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/utils"

	"github.com/jackc/pgx/v5"
)

type CartRepository interface {
	FindByUserId(ctx context.Context, tx pgx.Tx, userId string) (Cart, error)
	Save(ctx context.Context, tx pgx.Tx, cart Cart) Cart
	FindItemById(ctx context.Context, tx pgx.Tx, itemId string) (CartItem, error)
	FindItemByCartAndProduct(ctx context.Context, tx pgx.Tx, cartId string, productId string) (CartItem, error)
	SaveItem(ctx context.Context, tx pgx.Tx, cartItem CartItem) CartItem
	UpdateItem(ctx context.Context, tx pgx.Tx, cartItem CartItem) CartItem
	DeleteItem(ctx context.Context, tx pgx.Tx, itemId string)
	FindAllItems(ctx context.Context, tx pgx.Tx, cartId string) []CartItem
	DeleteAllItems(ctx context.Context, tx pgx.Tx, cartId string)
}

type CartRepositoryImpl struct{}

// DeleteAllItems implements [CartRepository].
func (c *CartRepositoryImpl) DeleteAllItems(ctx context.Context, tx pgx.Tx, cartId string) {
	query := "delete from cart_items where cart_id = $1"
	_, err := tx.Exec(ctx, query, cartId)
	utils.PanicIfError(err)
}

// DeleteItem implements [CartRepository].
func (c *CartRepositoryImpl) DeleteItem(ctx context.Context, tx pgx.Tx, itemId string) {
	query := "delete from cart_items where id = $1"
	_, err := tx.Exec(ctx, query, itemId)
	utils.PanicIfError(err)
}

// FindAllItems implements [CartRepository].
func (c *CartRepositoryImpl) FindAllItems(ctx context.Context, tx pgx.Tx, cartId string) []CartItem {
	query := "select id, cart_id, product_id, quantity, created_at from cart_items where cart_id = $1"
	cartItems := []CartItem{}
	rows, err := tx.Query(ctx, query, cartId)
	utils.PanicIfError(err)
	defer rows.Close()

	for rows.Next() {
		var cartItem CartItem
		rows.Scan(&cartItem.Id, &cartItem.CartId, &cartItem.ProductId, &cartItem.Quantity, &cartItem.CreatedAt)
		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		utils.PanicIfError(err)
	}

	return cartItems
}

// FindByUserId implements [CartRepository].
func (c *CartRepositoryImpl) FindByUserId(ctx context.Context, tx pgx.Tx, userId string) (Cart, error) {
	query := "select id, user_id, created_at, updated_at from carts where user_id = $1 limit 1"
	var cart Cart
	err := tx.QueryRow(ctx, query, userId).Scan(&cart.Id, &cart.UserId, &cart.CreatedAt, &cart.UpdatedAt)

	return cart, err
}

// FindItemByCartAndProduct implements [CartRepository].
func (c *CartRepositoryImpl) FindItemByCartAndProduct(ctx context.Context, tx pgx.Tx, cartId string, productId string) (CartItem, error) {
	query := "select id, cart_id, product_id, quantity, created_at from cart_items where cart_id = $1 and product_id = $2 limit 1"
	var cartItem CartItem
	err := tx.QueryRow(ctx, query, cartId, productId).Scan(&cartItem.Id, &cartItem.CartId, &cartItem.ProductId, &cartItem.Quantity, &cartItem.CreatedAt)
	return cartItem, err
}

// FindItemById implements [CartRepository].
func (c *CartRepositoryImpl) FindItemById(ctx context.Context, tx pgx.Tx, itemId string) (CartItem, error) {
	query := "select id, cart_id, product_id, quantity, created_at from cart_items where id = $1"
	var cartItem CartItem
	err := tx.QueryRow(ctx, query, itemId).Scan(&cartItem.Id, &cartItem.CartId, &cartItem.ProductId, &cartItem.Quantity, &cartItem.CreatedAt)
	return cartItem, err
}

// Save implements [CartRepository].
func (c *CartRepositoryImpl) Save(ctx context.Context, tx pgx.Tx, cart Cart) Cart {
	query := "insert into carts (user_id) values ($1) returning id, user_id, created_at, updated_at"
	var newCart Cart
	err := tx.QueryRow(ctx, query, cart.UserId).Scan(&newCart.Id, &newCart.UserId, &newCart.CreatedAt, &newCart.UpdatedAt)
	utils.PanicIfError(err)
	return newCart
}

// SaveItem implements [CartRepository].

func (c *CartRepositoryImpl) SaveItem(ctx context.Context, tx pgx.Tx, cartItem CartItem) CartItem {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		RETURNING id, cart_id, product_id, quantity, created_at
	`

	var newCartItem CartItem

	err := tx.QueryRow(
		ctx,
		query,
		cartItem.CartId,
		cartItem.ProductId,
		cartItem.Quantity,
	).Scan(
		&newCartItem.Id,
		&newCartItem.CartId,
		&newCartItem.ProductId,
		&newCartItem.Quantity,
		&newCartItem.CreatedAt,
	)
	utils.PanicIfError(err)

	return newCartItem
}

// UpdateItem implements [CartRepository].
func (c *CartRepositoryImpl) UpdateItem(ctx context.Context, tx pgx.Tx, cartItem CartItem) CartItem {
	query := `update cart_items set
		quantity = $1
	where id = $2 and cart_id = $3
	returning id, cart_id, product_id, quantity, created_at`

	var updatedCartItem CartItem
	err := tx.QueryRow(ctx, query, cartItem.Quantity, cartItem.Id, cartItem.CartId).Scan(&updatedCartItem.Id, &updatedCartItem.CartId, &updatedCartItem.ProductId, &updatedCartItem.Quantity, &updatedCartItem.CreatedAt)
	utils.PanicIfError(err)
	return updatedCartItem
}

func NewCartRepository() CartRepository {
	return &CartRepositoryImpl{}
}
