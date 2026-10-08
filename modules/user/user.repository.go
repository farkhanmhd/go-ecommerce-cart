package user

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	FindById(ctx context.Context, tx pgx.Tx, userId string) User
}

type UserRepositoryImpl struct{}

// FindById implements [UserRepository].
func (u *UserRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, userId string) User {
	query := "select id, name, email, created_at, updated_at from users where id = $1 limit 1"
	user := User{}
	err := tx.QueryRow(ctx, query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		panic(exceptions.NewNotFoundError("NOT FOUND"))
	}

	return user
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}
