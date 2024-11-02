package response

type GeneralResponse struct {
	StatusCode int    `json:"status_code"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

type DataResponse struct {
	GeneralResponse
	//make data flexible with any property
	Data interface{} `json:"data"`
}
