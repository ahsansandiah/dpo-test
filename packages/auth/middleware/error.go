package middlewareAuth

import "errors"

var (
	ErrorAuthHeaderEmpty       = errors.New("authorization header is empty")
	ErrorAuthNotHaveBearer     = errors.New("authorization header doesn't have bearer format")
	ErrorAuthNotHaveToken      = errors.New("authorization header doesn't have access token value")
	ErrorAccessTokenEmpty      = errors.New("failed getting data from access token")
	ErrorDataFromContext       = errors.New("failed getting data from context")
	ErrorInvalidTokenOrExpired = errors.New("token is invalid or has expired")
)
