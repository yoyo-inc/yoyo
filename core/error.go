package core

// ParameterError represents parameter error
type ParameterError struct {
	Message string
}

func (e ParameterError) Error() string {
	return e.Message
}

// NewParameterError returns parameter error
func NewParameterError(message string) ParameterError {
	return ParameterError{
		Message: message,
	}
}

// BusinessError represents bussiness error
type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

// NewBusinessError returns business error
func NewBusinessError(code string, message string) BusinessError {
	return BusinessError{
		Code:    code,
		Message: message,
	}
}
