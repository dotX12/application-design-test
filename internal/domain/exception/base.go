package exception

import (
	"errors"
)

var ErrDomain = errors.New("domain exception occurred")
var ErrApplication = errors.New("application exception occurred")
