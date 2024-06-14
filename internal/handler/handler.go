package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/umalmyha/oauth2-sample/internal/oauth2"
)

type Handler struct {
	creds  *OAuthCredentials
	iss    *oauth2.TokenIssuer
	logger *slog.Logger
}

func NewHandler(
	creds *OAuthCredentials,
	iss *oauth2.TokenIssuer,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		creds:  creds,
		iss:    iss,
		logger: logger,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id, secret, ok := r.BasicAuth()
	if !ok {
		h.logger.Error("Authorization header is missing or hs invalid format")
		h.respondJSON(w, http.StatusUnauthorized, InvalidClientError)
		return
	}

	if err := h.creds.Verify(id, secret); err != nil {
		h.logger.Error("failed to verify client id and secret", slog.Any("error", err))
		h.respondJSON(w, http.StatusUnauthorized, InvalidClientError)
		return
	}

	tkn, err := h.iss.Issue(time.Now().UTC())
	if err != nil {
		h.logger.Error("failed to issue JWT", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, http.StatusOK, tkn)
}

func (h *Handler) respondJSON(w http.ResponseWriter, code int, data any) {
	b, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("failed to marshal JSON", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err = w.Write(b); err != nil {
		h.logger.Error("failed to write JSON response", slog.Any("error", err))
		return
	}
}
