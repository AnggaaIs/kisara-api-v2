package response

type GeneralResponse struct {
	StatusCode int         `json:"status_code"`
	Name       string      `json:"name"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

func Success(statusCode int, name, message string, data interface{}) GeneralResponse {
	return GeneralResponse{
		StatusCode: statusCode,
		Name:       name,
		Message:    message,
		Data:       data,
	}
}

func ErrorWithDetails(statusCode int, name, message string, errors interface{}) GeneralResponse {
	return GeneralResponse{
		StatusCode: statusCode,
		Name:       name,
		Message:    message,
		Errors:     errors,
	}
}

func Error(statusCode int, name, message string, err error) GeneralResponse {
	var errorDetails interface{}
	if err != nil {
		errorDetails = map[string]interface{}{
			"details": err.Error(),
		}
	}

	return GeneralResponse{
		StatusCode: statusCode,
		Name:       name,
		Message:    message,
		Error:      errorDetails,
	}
}
