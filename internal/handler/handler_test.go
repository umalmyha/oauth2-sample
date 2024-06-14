package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/umalmyha/oauth2-sample/internal/handler"
	"github.com/umalmyha/oauth2-sample/internal/oauth2"
)

const (
	clientID     = "SLacIguEVzPheyKozLsa"
	clientSecret = "QjiqfOLwuOSryVRgqaFM"

	keyBits = 2048

	tokenIssuer = "oauth2-issuer"
	tokenTTL    = 10 * time.Minute

	endpoint = "http://localhost/token"
)

var privateKey, _ = rsa.GenerateKey(rand.Reader, keyBits)

func TestHandler_AccessTokenIssued(t *testing.T) {
	cut := handler.NewHandler(
		&handler.OAuthCredentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		oauth2.NewTokenIssuer(*privateKey, oauth2.WithIssuer(tokenIssuer), oauth2.WithTTL(tokenTTL)),
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelError,
				},
			),
		),
	)

	w := httptest.NewRecorder()
	r, err := request(clientID, clientSecret)
	if err != nil {
		t.Fatalf("failed construct request: %v", err)
	}

	cut.Handle(w, r)

	var tkn oauth2.AccessToken
	if err = json.NewDecoder(w.Body).Decode(&tkn); err != nil {
		t.Fatalf("failed decode response JSON body: %v", err)
	}

	if tkn.Type != oauth2.TokenTypeBearer {
		t.Fatalf("expected token type %v, received %v", oauth2.TokenTypeBearer, tkn.Type)
	}

	if seconds := int(tokenTTL.Seconds()); tkn.ExpiresIn != seconds {
		t.Fatalf("expected expires in %v, received %v", seconds, tkn.ExpiresIn)
	}

	var claims jwt.RegisteredClaims
	_, err = jwt.ParseWithClaims(tkn.AccessToken, &claims, func(token *jwt.Token) (any, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse sign token: %v", err)
	}

	if claims.Issuer != tokenIssuer {
		t.Fatalf("expected token issuer %v, received %v", tokenIssuer, claims.Issuer)
	}
}

func TestHandler_InvalidClient(t *testing.T) {
	cut := handler.NewHandler(
		&handler.OAuthCredentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		oauth2.NewTokenIssuer(*privateKey, oauth2.WithIssuer(tokenIssuer), oauth2.WithTTL(tokenTTL)),
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelError,
				},
			),
		),
	)

	w := httptest.NewRecorder()
	r, err := request("invalid_id", "invalid_secret")
	if err != nil {
		t.Fatalf("failed construct request: %v", err)
	}

	cut.Handle(w, r)

	var oauthErr handler.OAuthError
	if err = json.NewDecoder(w.Body).Decode(&oauthErr); err != nil {
		t.Fatalf("failed decode response JSON body: %v", err)
	}

	if oauthErr.Error != handler.InvalidClientError.Error {
		t.Fatalf("expected error %v, received %v", handler.InvalidClientError.Error, oauthErr.Error)
	}
}

func request(clientID, clientSecret string) (*http.Request, error) {
	r, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	return r, nil
}
