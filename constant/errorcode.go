package constant

const (
	SUCCESS        = 200
	INVALID_PARAMS = 400
	ERROR          = 500
)

const (
	ERROR_BAD_REQUEST = iota + 10000
	ERROR_DATABASE
	ERROR_USER_NOT_FOUND
	ERROR_USER_EXISTS
	ERROR_INCORRECT_PASSWORD
	ERROR_AUTH_CHECK_TOKEN_FAIL
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT
	ERROR_AUTH_CHECK_PERMISSION_DENIED
	ERROR_TOKEN_CLAIMS_PARSING_FAILED
	ERROR_TOKEN_CLAIMS_MISSING_FIELDS
	ERROR_FORBIDDEN
)

type ErrorCode struct {
	code int
}

func (e *ErrorCode) Error() string {
	return GetMsg(e.code)
}

func (e *ErrorCode) Code() int {
	return e.code
}

func NewError(code int) error {
	return &ErrorCode{code}
}
