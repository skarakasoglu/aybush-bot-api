package response

type Response struct{
	Code int `json:"code"`
	Error string `json:"error"`
}