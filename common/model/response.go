package common_model

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type ResponseWithData struct {
	Response
	Data any `json:"data"`
}
