package auth

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	Register(ctx context.Context, tx pgx.Tx, body RegisterUserBody)
	FindById(ctx context.Context, tx pgx.Tx, userId string) (User, error)
	FindByEmail(ctx context.Context, tx pgx.Tx, email string) (User, error)
	CreateSession(ctx context.Context, tx pgx.Tx, session Session)
	FindSessionByTokenHash(ctx context.Context, tx pgx.Tx, tokenHash string) Session
	RevokeSession(ctx context.Context, tx pgx.Tx, session Session)
	UpdateLastActive(ctx context.Context, tx pgx.Tx, session Session)
}

type AuthRepositoryImpl struct{}

// FindById implements [UserRepository].
func (u *AuthRepositoryImpl) FindById(ctx context.Context, tx pgx.Tx, userId string) (User, error) {
	query := "select id, name, email, password, created_at, updated_at from users where id = $1 limit 1"
	user := User{}
	err := tx.QueryRow(ctx, query, userId).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *AuthRepositoryImpl) FindByEmail(ctx context.Context, tx pgx.Tx, email string) (User, error) {
	query := "select id, name, email, password, created_at, updated_at from users where email = $1 limit 1"
	user := User{}
	err := tx.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return User{}, err
	} else {
		return user, nil
	}
}

// Save implements [UserRepository].
func (u *AuthRepositoryImpl) Register(ctx context.Context, tx pgx.Tx, body RegisterUserBody) {
	query := "insert into users (name, email, password) values ($1, $2, $3) returning id, name, email"
	var newUser User
	err := tx.QueryRow(ctx, query, body.Name, body.Email, body.Password).Scan(&newUser.Id, &newUser.Name, &newUser.Email)
	utils.PanicIfError(err)
}

func (u *AuthRepositoryImpl) CreateSession(ctx context.Context, tx pgx.Tx, session Session) {
	query := "insert into sessions (id, user_id, refresh_token_hash, expires_at, device, ip_address, user_agent) values ($1, $2, $3, $4, $5, $6, $7) returning id"
	var newSession Session
	err := tx.QueryRow(ctx, query, session.Id, session.UserId, session.RefreshTokenHash, session.ExpiresAt, session.Device, session.IpAdress, session.UserAgent).Scan(&newSession.Id)
	utils.PanicIfError(err)
}

func (u *AuthRepositoryImpl) FindSessionByTokenHash(ctx context.Context, tx pgx.Tx, tokenHash string) Session {
	query := "select id, user_id, expires_at from sessions where refresh_token_hash = $1 limit 1"
	session := Session{}
	err := tx.QueryRow(ctx, query, tokenHash).Scan(&session.Id, &session.UserId, &session.ExpiresAt)

	if err != nil {
		panic(exceptions.NewAuthorizationError("Invalid or Expired Token"))
	}

	return session
}

func (u *AuthRepositoryImpl) RevokeSession(ctx context.Context, tx pgx.Tx, session Session) {
	query := "delete from sessions where id = $1 and user_id = $2"
	_, err := tx.Exec(ctx, query, session.Id, session.UserId)

	if err != nil {
		panic(err)
	}
}

func (u *AuthRepositoryImpl) UpdateLastActive(ctx context.Context, tx pgx.Tx, session Session) {
	query := "update sessions set last_active_at = $1 where id = $2"
	_, err := tx.Exec(ctx, query, time.Now(), session.Id)
	utils.PanicIfError(err)
}

func NewAuthRepository() AuthRepository {
	return &AuthRepositoryImpl{}
}
