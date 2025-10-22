package domain

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
	Response struct {
		Message    string `json:"message"`
		StatusCode int    `json:"statuscode"`
	}
)
