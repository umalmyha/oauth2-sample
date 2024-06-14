package oauth2

import "time"

const (
	DefaultTokenIssuer = "oauth2-server-issuer"
	DefaultTokenTTL    = 5 * time.Minute
)

type ConfigFunc func(cfg *Config)

type Config struct {
	issuer     string
	ttlDur     time.Duration
	ttlSeconds int
}

func WithTTL(ttl time.Duration) ConfigFunc {
	return func(cfg *Config) {
		if ttl > 0 {
			cfg.ttlDur = ttl
			cfg.ttlSeconds = int(ttl.Seconds())
		}
	}
}
func WithIssuer(iss string) ConfigFunc {
	return func(cfg *Config) {
		if iss != "" {
			cfg.issuer = iss
		}
	}
}
