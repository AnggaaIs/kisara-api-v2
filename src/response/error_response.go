package response

//create extendable response from GeneralResponse
type NotFoundResponse struct {
	GeneralResponse
	Path string `json:"path"`
}

type ErrorField struct {
	Name    string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	GeneralResponse
	Errors []ErrorField `json:"errors"`
}
