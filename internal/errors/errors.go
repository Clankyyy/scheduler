package errors

type APIError struct {
	HTTPStatusCode int    `json:"httpStatusCode"`
	Message        string `json:"message"`
	Details        string `json:"details,omitempty"`
}

func NewAPIError(status uint, msg string) APIError {
	return APIError{
		HTTPStatusCode: int(status),
		Message:        msg,
	}
}
