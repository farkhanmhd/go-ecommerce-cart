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
	GetCart(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	AddItem(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdateItem(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	RemoveItem(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	ClearCart(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type CartHandlerImpl struct {
	CartService CartService
}

func (c *CartHandlerImpl) ClearCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}
	c.CartService.ClearCart(r.Context(), claims.UserId)
	utils.WriteJSONResponse(w, common_model.Response{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
	})
}

func (c *CartHandlerImpl) RemoveItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}
	itemId := p.ByName("itemId")
	c.CartService.RemoveItem(r.Context(), claims.UserId, itemId)
	utils.WriteJSONResponse(w, common_model.Response{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
	})
}

func (c *CartHandlerImpl) UpdateItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}

	itemId := p.ByName("itemId")

	body := CartItem{}
	utils.ReadFromRequestBody(r, &body)
	body.Id = itemId

	cartItem := c.CartService.UpdateItem(r.Context(), claims.UserId, body)
	response := common_model.ResponseWithData{
		Data: cartItem,
	}
	response.Code = http.StatusOK
	response.Status = http.StatusText(http.StatusOK)
	utils.WriteJSONResponse(w, response)
}

func (c *CartHandlerImpl) GetCart(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		panic(exceptions.NewAuthorizationError("UNAUTHORIZED"))
	}

	cartItems := c.CartService.GetCart(r.Context(), claims.UserId)
	response := common_model.ResponseWithData{
		Data: cartItems,
	}
	response.Code = http.StatusOK
	response.Status = http.StatusText(http.StatusOK)
	utils.WriteJSONResponse(w, response)
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
