package exceptions

import (
	common_model "farkhanmhd/go-ecommerce-cart/common/model"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func internalServerError(w http.ResponseWriter, _ *http.Request, _ any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	response := common_model.Response{
		Code:   http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
	}

	utils.WriteJSONResponse(w, response)
}

func validationError(w http.ResponseWriter, _ *http.Request, err any) bool {
	exception, ok := err.(validator.ValidationErrors)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := common_model.ResponseWithData{
			Data: exception.Error(),
		}

		response.Code = http.StatusBadRequest
		response.Status = "BAD REQUEST"

		utils.WriteJSONResponse(w, response)
		return true
	}

	return false
}

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

func notFoundError(w http.ResponseWriter, _ *http.Request, err any) bool {
	exception, ok := err.(NotFoundError)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := common_model.ResponseWithData{
			Data: exception.Error,
		}

		response.Code = http.StatusNotFound
		response.Status = "NOT FOUND"

		utils.WriteJSONResponse(w, response)
		return true
	}

	return false
}

type AuthorizationError struct {
	Error string
}

func NewAuthorizationError(error string) AuthorizationError {
	return AuthorizationError{Error: error}
}

func authorizationError(w http.ResponseWriter, _ *http.Request, err any) bool {
	exception, ok := err.(AuthorizationError)

	if ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		response := common_model.ResponseWithData{
			Data: exception.Error,
		}

		response.Code = http.StatusUnauthorized
		response.Status = "UNAUTHORIZED"

		utils.WriteJSONResponse(w, response)
		return true
	}

	return false
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err any) {
	if validationError(w, r, err) {
		return
	}

	if notFoundError(w, r, err) {
		return
	}

	if authorizationError(w, r, err) {
		return
	}

	internalServerError(w, r, err)
}
