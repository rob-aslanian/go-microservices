package hc_errors

const (
	GENERAL_ERROR_TYPE            = "GeneralError"
	MONGO_ERROR_TYPE              = "MongoError"
	DUPLICATE_KEY_ERROR_TYPE      = "DuplicateKeyError"
	AUTH_ERROR_TYPE               = "AuthError"
	NOT_FOUNT_ERROR_TYPE          = "NotFountError"
	PARAMETER_MISMATCH_ERROR_TYPE = "ParameterMismatchError"
)

var NOT_AUTHENTICATED_ERROR = JsonError{
	Type:        AUTH_ERROR_TYPE,
	Description: "You are not authenticated",
}

var INVALID_LOGIN_PASSWORD_ERROR = JsonError{
	Type:        AUTH_ERROR_TYPE,
	Description: "Invalid Login or Password",
}

var USER_NOT_FOUND = JsonError{
	Type:        NOT_FOUNT_ERROR_TYPE,
	Description: "User not found",
}

var TOKEN_NOT_FOUND = JsonError{
	Type:        NOT_FOUNT_ERROR_TYPE,
	Description: "Token not found",
}

var TOKEN_DOES_NOT_MATCH = JsonError{
	Type:        PARAMETER_MISMATCH_ERROR_TYPE,
	Description: "Token does not match user",
}

var DUPLICATE_KEY_ERROR = JsonError{
	Type:        DUPLICATE_KEY_ERROR_TYPE,
	Description: "The value you entered already exists",
}

var USER_IS_ALREADY_ACTIVATED = JsonError{
	Type:        GENERAL_ERROR_TYPE,
	Description: "User is already activated",
}
