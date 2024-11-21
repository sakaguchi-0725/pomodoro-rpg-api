package apperr

import "github.com/cockroachdb/errors"

var (
	ErrDataNotFound        = errors.New("DataNotFound")
	ErrInvalidParameter    = errors.New("InvalidParameter")
	ErrUnautorizedExeption = errors.New("Unauthorized")
)
