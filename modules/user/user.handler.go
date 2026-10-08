package user

import (
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	common_model "farkhanmhd/go-ecommerce-cart/common/model"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserHandler interface {
	FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type UserHandlerImpl struct {
	UserService UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &UserHandlerImpl{
		UserService: userService,
	}
}

// FindAll implements [UserHandler].
func (u *UserHandlerImpl) FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())

	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}

	// claims.UserId valid. error on service
	usersJson := u.UserService.FindById(r.Context(), claims.UserId)
	response := common_model.ResponseWithData{
		Data: usersJson,
	}

	response.Code = http.StatusOK
	response.Status = "OK"

	utils.WriteJSONResponse(w, response)
}
