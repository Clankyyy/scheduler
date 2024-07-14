package errors

type APIError struct {
	StatusCode int    `json:"httpStatusCode"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

func (ae APIError) Error() string {
	return ae.Message
}

func NewAPIError(status uint, msg string) APIError {
	return APIError{
		StatusCode: int(status),
		Message:    msg,
	}
}
