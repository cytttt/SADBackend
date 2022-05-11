package constant

var MsgFlags = map[int]string{

	SUCCESS:                            "Ok",
	INVALID_PARAMS:                     "Invalid params error",
	ERROR:                              "Fail",
	ERROR_USER_NOT_FOUND:               "User not exist or incorrect password",
	ERROR_AUTH_CHECK_TOKEN_FAIL:        "Auth token check fail",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:     "Auth token expired",
	ERROR_AUTH_CHECK_PERMISSION_DENIED: "Permission denied",
	ERROR_BAD_REQUEST:                  "Bad request",
	ERROR_TOKEN_CLAIMS_PARSING_FAILED:  "Failed to parse token claims",
	ERROR_TOKEN_CLAIMS_MISSING_FIELDS:  "Some fields are missing in token claims",
	ERROR_FORBIDDEN:                    "Response forbidden",
	ERROR_DUPLICATED_PIXEL_NAME:        "Duplicated pixel name",
	ERROR_CDPAPI_BAD_REQUEST:           "CDP API bad request",
	ERROR_ECSERVER_BAD_REQUEST:         "ECServer bad request",
	ERROR_INVALID_CSV_CONTENT:          "Invalid csv content",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
