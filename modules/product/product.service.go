package product

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/utils"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	FindById(ctx context.Context, ProductId string) Product
	Create(ctx context.Context, product Product) Product
	FindAll(ctx context.Context) []Product
	Update(ctx context.Context, product Product)
	Delete(ctx context.Context, productId string)
}

type ProductServiceImpl struct {
	ProductRepository ProductRepository
	Validate          *validator.Validate
	Pool              *pgxpool.Pool
}

func NewProductService(ProductRepository ProductRepository, pool *pgxpool.Pool, validate *validator.Validate) ProductService {
	return &ProductServiceImpl{
		ProductRepository: ProductRepository,
		Pool:              pool,
		Validate:          validate,
	}
}

func (p *ProductServiceImpl) FindById(ctx context.Context, ProductId string) Product {
	tx, err := p.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	product := p.ProductRepository.FindById(ctx, tx, ProductId)
	return product
}

func (p *ProductServiceImpl) Create(ctx context.Context, product Product) Product {
	err := p.Validate.Struct(product)
	utils.PanicIfError(err)
	tx, err := p.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	newProduct := p.ProductRepository.Create(ctx, tx, product)
	return newProduct
}

func (p *ProductServiceImpl) FindAll(ctx context.Context) []Product {
	tx, err := p.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	products := p.ProductRepository.FindAll(ctx, tx)
	return products
}

func (p *ProductServiceImpl) Update(ctx context.Context, product Product) {
	err := p.Validate.Struct(product)
	utils.PanicIfError(err)
	tx, err := p.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	p.ProductRepository.FindById(ctx, tx, product.Id)
	p.ProductRepository.Update(ctx, tx, product)
}

func (p *ProductServiceImpl) Delete(ctx context.Context, productId string) {
	tx, err := p.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	p.ProductRepository.FindById(ctx, tx, productId)
	p.ProductRepository.Delete(ctx, tx, productId)
}
