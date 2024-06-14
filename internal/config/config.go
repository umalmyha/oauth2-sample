package config

import (
	"crypto/rsa"
	"reflect"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	HTTPServer  HTTPServer  `envPrefix:"HTTP_SERVER_"`
	JWT         JWT         `envPrefix:"JWT_"`
	Credentials Credentials `envPrefix:"OAUTH_CREDENTIALS_"`
}

type JWT struct {
	PrivateKey rsa.PrivateKey `env:"PRIVATE_KEY,notEmpty"`
	Issuer     string         `env:"ISSUER"`
	TTL        time.Duration  `env:"TTL"`
}

type Credentials struct {
	ClientID     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
}

type HTTPServer struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func Load() (cfg Config, err error) {
	fm := map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(rsa.PrivateKey{}): func(v string) (any, error) {
			pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(v))
			if err != nil {
				return nil, err
			}
			return *pk, nil
		},
	}

	err = env.ParseWithOptions(&cfg, env.Options{FuncMap: fm})

	return cfg, err
}
