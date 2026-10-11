package user

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService interface {
	FindById(ctx context.Context, userId string) UserResponse
}

type UserServiceImpl struct {
	UserRepository UserRepository
	Pool           *pgxpool.Pool
}

func NewUserService(userRepository UserRepository, pool *pgxpool.Pool) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		Pool:           pool,
	}
}

func (u *UserServiceImpl) FindById(ctx context.Context, userId string) UserResponse {
	tx, err := u.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)
	user := u.UserRepository.FindById(ctx, tx, userId)
	response := UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return response
}
