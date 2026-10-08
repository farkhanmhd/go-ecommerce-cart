package token

import (
	"crypto/sha256"
	"encoding/hex"
	"farkhanmhd/go-ecommerce-cart/common/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

type Claims struct {
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type TokenSigner struct {
	key      []byte
	issuer   string
	audience string
	ttl      time.Duration
}

func (t *TokenSigner) Sign(u entity.User, sessionId string) (string, error) {
	now := time.Now()
	claims := Claims{
		Name:      u.Name,
		UserId:    u.Id,
		SessionId: sessionId,
		Email:     u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   u.Id,
			Audience:  jwt.ClaimStrings{t.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.key)
}

func (t *TokenSigner) Validate(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(*jwt.Token) (any, error) {
			return t.key, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(t.issuer),
		jwt.WithAudience(t.audience),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func NewTokenSign(key []byte, issuer, audience string, ttl time.Duration) *TokenSigner {
	return &TokenSigner{
		key:      key,
		issuer:   issuer,
		audience: audience,
		ttl:      ttl,
	}
}
