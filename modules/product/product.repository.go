package product

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	FindById(ctx context.Context, tx pgx.Tx, productId string) Product
	Create(ctx context.Context, tx pgx.Tx, product Product) Product
	FindAll(ctx context.Context, tx pgx.Tx) []Product
	Update(ctx context.Context, tx pgx.Tx, product Product)
	Delete(ctx context.Context, tx pgx.Tx, productId string)
}

type ProductRepositoryImpl struct{}

func (p *ProductRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, productId string) Product {
	query := "select id, name, description, price, stock, created_at, updated_at from products where id = $1 limit 1"
	product := Product{}
	err := tx.QueryRow(ctx, query, productId).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		panic(exceptions.NewNotFoundError("NOT FOUND"))
	}

	return product
}

func (p *ProductRepositoryImpl) Create(ctx context.Context, tx pgx.Tx, product Product) Product {
	query := "insert into products (id, name, description, price, stock) values ($1, $2, $3, $4, $5) returning id, name, description, price, stock, created_at, updated_at"
	newProductId := uuid.New()
	newProduct := Product{}
	err := tx.QueryRow(ctx, query, newProductId, product.Name, product.Description, product.Price, product.Stock).Scan(
		&newProduct.Id,
		&newProduct.Name,
		&newProduct.Description,
		&newProduct.Price,
		&newProduct.Stock,
		&newProduct.CreatedAt,
		&newProduct.UpdatedAt)

	utils.PanicIfError(err)
	return newProduct
}

func (p *ProductRepositoryImpl) FindAll(ctx context.Context, tx pgx.Tx) []Product {
	query := "select id, name, description, price, stock, created_at, updated_at from products"
	products := []Product{}
	rows, err := tx.Query(ctx, query)
	utils.PanicIfError(err)
	defer rows.Close()

	for rows.Next() {
		var productRow Product
		rows.Scan(&productRow.Id, &productRow.Name, &productRow.Description, &productRow.Price, &productRow.Stock, &productRow.CreatedAt, &productRow.UpdatedAt)
		products = append(products, productRow)
	}

	if err := rows.Err(); err != nil {
		utils.PanicIfError(err)
	}

	return products
}

func (p *ProductRepositoryImpl) Update(ctx context.Context, tx pgx.Tx, product Product) {
	query := `update products set
		name = $1,
		description = $2,
		price = $3,
	    stock = $4,
		updated_at = $5
	 where id = $6`
	_, err := tx.Exec(ctx, query, product.Name, product.Description, product.Price, product.Stock, time.Now(), product.Id)
	utils.PanicIfError(err)
}

func (p *ProductRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, productId string) {
	query := "delete from products where id = $1"
	_, err := tx.Exec(ctx, query, productId)
	utils.PanicIfError(err)
}

func NewProductRepository() ProductRepository {
	return &ProductRepositoryImpl{}
}
