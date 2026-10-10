package auth

import (
	"farkhanmhd/go-ecommerce-cart/common/entity"
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	common_model "farkhanmhd/go-ecommerce-cart/common/model"
	"farkhanmhd/go-ecommerce-cart/common/token"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Login(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Refresh(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type AuthHandlerImpl struct {
	AuthService  AuthService
	Signer       *token.TokenSigner
	CookieSecure bool
}

func NewAuthHandler(authService AuthService, signer *token.TokenSigner) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
		Signer:      signer,
	}
}

const REFRESH_COOKIE_NAME = "GO_ECOMMERCE_CART_REFRESH_TOKEN"

func (u *AuthHandlerImpl) setRefreshCookie(w http.ResponseWriter, raw string) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_COOKIE_NAME,
		Value:    raw,
		Path:     "/api",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   u.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func (u *AuthHandlerImpl) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_COOKIE_NAME,
		Value:    "",
		MaxAge:   -1,
		Path:     "/api",
		HttpOnly: true,
		Secure:   u.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func (u *AuthHandlerImpl) Register(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	createUserBody := RegisterUserBody{}
	utils.ReadFromRequestBody(r, &createUserBody)

	u.AuthService.Register(r.Context(), createUserBody)
	response := common_model.Response{
		Code:   http.StatusCreated,
		Status: "CREATED",
	}

	utils.WriteJSONResponse(w, response)
}

func (u *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	loginBody := LoginBody{}
	utils.ReadFromRequestBody(r, &loginBody)
	userAgent := r.UserAgent()

	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	sessionId := uuid.New().String()

	result := u.AuthService.Login(r.Context(), loginBody, userAgent, ipAddress, sessionId)

	accessToken, err := u.Signer.Sign(entity.User{
		Id:        result.User.Id,
		Name:      result.User.Name,
		Email:     result.User.Email,
		CreatedAt: result.User.CreatedAt,
		UpdatedAt: result.User.UpdatedAt,
	}, sessionId)

	utils.PanicIfError(err)

	u.setRefreshCookie(w, result.RefreshToken)
	utils.WriteJSONResponse(w, common_model.ResponseWithData{
		Response: common_model.Response{
			Status: "OK",
			Code:   http.StatusOK,
		},
		Data: LoginResponse{
			AccessToken: accessToken,
			User: LoginUserResult{
				Id:        result.User.Id,
				Name:      result.User.Name,
				Email:     result.User.Email,
				CreatedAt: result.User.CreatedAt,
				UpdatedAt: result.User.UpdatedAt,
			},
		},
	})
}

func (u *AuthHandlerImpl) Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cookie, err := r.Cookie(REFRESH_COOKIE_NAME)
	if err == nil && cookie != nil {
		u.AuthService.Logout(r.Context(), cookie.Value)
	}
	u.clearRefreshCookie(w)
	utils.WriteJSONResponse(w, common_model.Response{
		Status: "OK",
		Code:   http.StatusOK,
	})
}

func (u *AuthHandlerImpl) Refresh(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cookie, err := r.Cookie(REFRESH_COOKIE_NAME)
	if err != nil {
		panic(exceptions.NewAuthorizationError(exceptions.SessionExpiredMessage))
	}
	user := u.AuthService.Refresh(r.Context(), cookie.Value)

	accessToken, err := u.Signer.Sign(entity.User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, cookie.Value)

	utils.WriteJSONResponse(w, common_model.ResponseWithData{
		Response: common_model.Response{
			Status: "OK",
			Code:   http.StatusOK,
		},
		Data: LoginResponse{
			AccessToken: accessToken,
			User: LoginUserResult{
				Id:        user.Id,
				Name:      user.Name,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
	})
}
