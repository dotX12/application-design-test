package exception

import "fmt"

var ErrInvalidDateRange = fmt.Errorf("%w: invalid date range", ErrDomain)
var ErrQuotaIsNegative = fmt.Errorf("%w: quouta is negative", ErrDomain)
var ErrEmailFormat = fmt.Errorf("%w: invalid email format", ErrDomain)
