package customerror

type APIError struct {
	Code    int
	Message string
}

func NewAPIError(message string, code int) APIError {
	return APIError{
		Message: message,
		Code:    code,
	}
}
