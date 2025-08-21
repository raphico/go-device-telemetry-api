package http

type errorCode string

const (
	invalidRequest errorCode = "INVALID_REQUEST" // includes validation errors
	unauthorized   errorCode = "UNAUTHORIZED"    // invalid token, creds, grant
	conflict       errorCode = "CONFLICT"        // duplicate email/username
	internalError  errorCode = "INTERNAL_ERROR"  // unexpected server error
	notfound       errorCode = "NOT_FOUND"
)
