package apperr

var (
	ErrBadRequest          = "request is incomplete or invalid"
	ErrUnprocessable       = "request is complete but invalid "
	ErrInvalidParam        = "invalid input param(s)"
	ErrUnauthorized        = "authorization missing or bad"
	ErrPermissionDenied    = "permission denied"
	ErrNotFound            = "record not found"
	ErrInternalServerError = "internal server error"
	ErrDuplicateRequest    = "duplicate request"
	ErrTimeout             = "timeout"
	ErrConnectTimeout      = "connect timeout"
	ErrRateLimited         = "rate limited"
	ErrUserError           = "user error"
	ErrUnexpectedError     = "unexpected error"
	ErrTransientError      = "transient application error"
)

// User
var (
	ErrDuplicateUser = "user already exists"
)
