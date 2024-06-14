package oauth2

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenTypeBearer = "bearer"

type TokenIssuer struct {
	key *rsa.PrivateKey
	cfg *Config
}

func NewTokenIssuer(key rsa.PrivateKey, opts ...ConfigFunc) *TokenIssuer {
	cfg := Config{
		issuer:     DefaultTokenIssuer,
		ttlDur:     DefaultTokenTTL,
		ttlSeconds: int(DefaultTokenTTL.Seconds()),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	return &TokenIssuer{
		key: &key,
		cfg: &cfg,
	}
}

func (i *TokenIssuer) Issue(issueAt time.Time) (AccessToken, error) {
	claims := jwt.RegisteredClaims{
		Issuer:   i.cfg.issuer,
		IssuedAt: jwt.NewNumericDate(issueAt),
		ExpiresAt: jwt.NewNumericDate(
			issueAt.Add(i.cfg.ttlDur),
		),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(i.key)
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		AccessToken: signed,
		Type:        TokenTypeBearer,
		ExpiresIn:   i.cfg.ttlSeconds,
	}, nil
}
