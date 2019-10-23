package api

type Response struct {
	Ok bool `json:"ok"`
	Data interface{} `json:"data"`
	Error error `json:"error"`
}