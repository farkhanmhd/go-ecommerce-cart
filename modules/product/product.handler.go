package product

import (
	common_model "farkhanmhd/go-ecommerce-cart/common/model"
	"farkhanmhd/go-ecommerce-cart/common/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ProductHandler interface {
	FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Create(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	FindAll(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Update(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type ProductHandlerImpl struct {
	ProductService ProductService
}

func NewProductHandler(ProductService ProductService) ProductHandler {
	return &ProductHandlerImpl{
		ProductService: ProductService,
	}
}

// FindAll implements [ProductHandler].
func (i *ProductHandlerImpl) FindById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productId := p.ByName("productId")
	product := i.ProductService.FindById(r.Context(), productId)
	response := common_model.ResponseWithData{
		Data: product,
	}

	response.Code = http.StatusOK
	response.Status = "OK"

	utils.WriteJSONResponse(w, response)
}

func (i *ProductHandlerImpl) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	product := Product{}
	utils.ReadFromRequestBody(r, &product)

	i.ProductService.Create(r.Context(), product)
	response := common_model.Response{
		Code:   http.StatusCreated,
		Status: "CREATED",
	}

	utils.WriteJSONResponse(w, response)
}

func (i *ProductHandlerImpl) FindAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	products := i.ProductService.FindAll(r.Context())
	response := common_model.ResponseWithData{
		Data: products,
	}

	response.Code = http.StatusOK
	response.Status = "OK"

	utils.WriteJSONResponse(w, response)
}

func (i *ProductHandlerImpl) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productId := p.ByName("productId")
	product := Product{}
	utils.ReadFromRequestBody(r, &product)
	product.Id = productId

	i.ProductService.Update(r.Context(), product)
	response := common_model.Response{
		Code:   http.StatusOK,
		Status: "OK",
	}
	utils.WriteJSONResponse(w, response)
}

func (i *ProductHandlerImpl) Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	productId := p.ByName("productId")
	i.ProductService.Delete(r.Context(), productId)
	response := common_model.Response{
		Code:   http.StatusOK,
		Status: "OK",
	}
	utils.WriteJSONResponse(w, response)
}
