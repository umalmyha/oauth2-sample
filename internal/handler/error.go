package handler

// InvalidClientError https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
var InvalidClientError = OAuthError{Error: "invalid_client"}

type OAuthError struct {
	Error string `json:"error"`
}
