package http

type errorCode string

const (
	invalidRequest  errorCode = "INVALID_REQUEST"
	invalidEmail    errorCode = "INVALID_EMAIL"
	invalidUsername errorCode = "INVALID_USERNAME"
	invalidPassword errorCode = "INVALID_PASSWORD"
	unauthorized    errorCode = "UNAUTHORIZED"

	invalidCredentials errorCode = "INVALID_CREDENTIALS"
	invalidGrant       errorCode = "INVALID_GRANT"
	validationError    errorCode = "VALIDATION_ERROR"

	emailExists    errorCode = "EMAIL_EXISTS"
	usernameExists errorCode = "USERNAME_EXISTS"

	internalError errorCode = "INTERNAL_ERROR"
)
