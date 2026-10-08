package auth

import (
	"farkhanmhd/go-ecommerce-cart/common/entity"
	"time"
)

type Session struct {
	Id               string
	UserId           string
	RefreshTokenHash string
	ExpiresAt        time.Time
	CreatedAt        time.Time
	Device           string
	LastActiveAt     time.Time
	IpAdress         string
	UserAgent        string
}

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginUserResult struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUserBody struct {
	Name     string `validate:"required,max=200,min=1" json:"name"`
	Email    string `validate:"required,max=200,min=1" json:"email"`
	Password string `validate:"required,max=200,min=1" json:"password"`
}

type RegisterUserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func RegisterUserToJson(user entity.User) RegisterUserResponse {
	return RegisterUserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}

type LoginBody struct {
	Email    string `validate:"required,max=200,min=1" json:"email"`
	Password string `validate:"required,max=200,min=1" json:"password"`
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	User         User
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        LoginUserResult
}
