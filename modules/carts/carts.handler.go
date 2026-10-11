package carts

import (
	"farkhanmhd/go-ecommerce-cart/common/exceptions"
	"farkhanmhd/go-ecommerce-cart/common/middleware"
	common_model "farkhanmhd/go-ecommerce-cart/common/model"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CartHandler interface {
	AddItem(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type CartHandlerImpl struct {
	CartService CartService
}

func (c *CartHandlerImpl) AddItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}

	body := CartItem{}
	utils.ReadFromRequestBody(r, &body)

	result := c.CartService.AddItem(r.Context(), claims.UserId, body)
	response := common_model.ResponseWithData{
		Data: result,
	}

	response.Code = http.StatusCreated
	response.Status = http.StatusText(http.StatusCreated)

	utils.WriteJSONResponse(w, response)
}

func NewCartHandler(cartService CartService) CartHandler {
	return &CartHandlerImpl{
		CartService: cartService,
	}
}
