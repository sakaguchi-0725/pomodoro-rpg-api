package apperr

type ApplicationError struct {
	code    ErrorCode
	message string
	err     error
}

func (e *ApplicationError) Error() string {
	return e.err.Error()
}

func (e *ApplicationError) Code() ErrorCode {
	return e.code
}

func (e *ApplicationError) Message() string {
	return e.message
}

func NewApplicationError(code ErrorCode, message string, err error) *ApplicationError {
	return &ApplicationError{
		code:    code,
		message: message,
		err:     err,
	}
}
