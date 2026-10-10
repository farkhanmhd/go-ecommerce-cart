package auth

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/token"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService interface {
	Register(ctx context.Context, requestBody RegisterUserBody)
	Login(ctx context.Context, requestBody LoginBody, userAgent string, ipAddress string, sessionId string) LoginServiceResult
	Logout(ctx context.Context, refreshToken string)
	Refresh(ctx context.Context, refreshToken string) User
}

type AuthServiceImpl struct {
	AuthRepository AuthRepository
	Pool           *pgxpool.Pool
	Validate       *validator.Validate
	Signer         *token.TokenSigner
}

func NewAuthService(authRepository AuthRepository, pool *pgxpool.Pool, validate *validator.Validate, signer *token.TokenSigner) AuthService {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
		Pool:           pool,
		Validate:       validate,
		Signer:         signer,
	}
}

func (u *AuthServiceImpl) Register(ctx context.Context, requestBody RegisterUserBody) {
	err := u.Validate.Struct(requestBody)
	utils.PanicIfError(err)
	hashedPassword := utils.HashPassword(requestBody.Password)

	body := RegisterUserBody{
		Name:     requestBody.Name,
		Email:    requestBody.Email,
		Password: hashedPassword,
	}

	tx, err := u.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	u.AuthRepository.Register(ctx, tx, body)
}

func (u *AuthServiceImpl) Login(ctx context.Context, requestBody LoginBody, userAgent string, ipAddress string, sessionId string) LoginServiceResult {
	err := u.Validate.Struct(requestBody)
	utils.PanicIfError(err)

	tx, err := u.Pool.Begin(ctx)
	utils.PanicIfError(err)
	defer utils.CommitOrRollback(ctx, tx)

	userResult, err := u.AuthRepository.FindByEmail(ctx, tx, requestBody.Email)

	if err != nil {
		panic(exceptions.NewAuthorizationError("Wrong email or password"))
	}

	err = utils.ComparePassword(userResult.Password, requestBody.Password)

	if err != nil {
		panic(exceptions.NewAuthorizationError("Wrong email or password"))
	}

	rawRefreshToken := uuid.New().String()

	hashToken := token.HashToken(rawRefreshToken)
	device := userAgent

	if len(device) > 100 {
		device = device[:100]
	}

	session := Session{
		Id:               sessionId,
		UserId:           userResult.Id,
		RefreshTokenHash: hashToken,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
		Device:           device,
		IpAdress:         ipAddress,
		UserAgent:        userAgent,
	}

	u.AuthRepository.CreateSession(ctx, tx, session)

	return LoginServiceResult{
		RefreshToken: rawRefreshToken,
		User:         userResult,
	}
}

func (u *AuthServiceImpl) Logout(ctx context.Context, refreshToken string) {
	if len(refreshToken) == 0 {
		return
	}
	tx, err := u.Pool.Begin(ctx)
	defer utils.CommitOrRollback(ctx, tx)

	if err != nil {
		panic(err)
	}

	tokenHash := token.HashToken(refreshToken)
	session := u.AuthRepository.FindSessionByTokenHash(ctx, tx, tokenHash)
	u.AuthRepository.RevokeSession(ctx, tx, session)
}

func (u *AuthServiceImpl) Refresh(ctx context.Context, refreshToken string) User {
	if len(refreshToken) == 0 {
		panic(exceptions.NewAuthorizationError("Invalid Token"))
	}

	tx, err := u.Pool.Begin(ctx)
	defer utils.CommitOrRollback(ctx, tx)

	if err != nil {
		panic(err)
	}

	tokenHash := token.HashToken(refreshToken)
	session := u.AuthRepository.FindSessionByTokenHash(ctx, tx, tokenHash)

	if session.ExpiresAt.Before(time.Now()) {
		u.AuthRepository.RevokeSession(ctx, tx, session)
		panic(exceptions.NewAuthorizationError("Session Expired"))
	}

	u.AuthRepository.UpdateLastActive(ctx, tx, session)

	user, err := u.AuthRepository.FindById(ctx, tx, session.UserId)

	if err != nil {
		fmt.Println("Error here at find by id")
		panic(exceptions.NewAuthorizationError(exceptions.SessionExpiredMessage))
	}

	return user
}
