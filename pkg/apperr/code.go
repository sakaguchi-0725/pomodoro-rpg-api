package apperr

type ErrorCode int

const (
	ErrBadRequest ErrorCode = iota
	ErrNotFound
	ErrUnautorized
)

func (c ErrorCode) String() string {
	switch c {
	case ErrBadRequest:
		return "BadRequest"
	case ErrNotFound:
		return "NotFound"
	case ErrUnautorized:
		return "Unautorized"
	default:
		return "InternalServerError"
	}
}
