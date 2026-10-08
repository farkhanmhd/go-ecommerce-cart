package middleware

import (
	"context"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/token"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type UserContextKey string

const UserClaimsKey UserContextKey = "user_claims"

type AuthMiddleware struct {
	signer *token.TokenSigner
}

func NewAuthMiddleware(signer *token.TokenSigner) *AuthMiddleware {
	return &AuthMiddleware{signer: signer}
}

func (m *AuthMiddleware) Protect(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			panic(exceptions.NewAuthorizationError("Missing or invalid authorization header"))
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.signer.Validate(rawToken)
		if err != nil {
			panic(exceptions.NewAuthorizationError("Unauthorized: " + err.Error()))
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next(w, r.WithContext(ctx), p)
	}
}

func GetUserClaims(ctx context.Context) (*token.Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*token.Claims)
	return claims, ok
}
