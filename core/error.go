package core

type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewBusinessError returns runtime error
func NewBusinessError(code string, message string) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}
