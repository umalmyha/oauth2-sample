package handler

import "errors"

type OAuthCredentials struct {
	ClientID     string
	ClientSecret string
}

func (c *OAuthCredentials) Verify(id, secret string) error {
	if c.ClientID != id || c.ClientSecret != secret {
		return errors.New("incorrect client credentials")
	}
	return nil
}
