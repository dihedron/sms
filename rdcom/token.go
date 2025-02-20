package rdcom

import (
	"errors"
	"log/slog"
	"time"

	"github.com/dihedron/sms/pointer"
)

type TokenService struct {
	Service
}

type Token struct {
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expire_date"`
}

// List returns the list of tokens.
func (t *TokenService) List() ([]Token, error) {
	if t.backref.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}
	tokens, err := doGet[Token](t.backref, &GetRequest{
		Request: Request{
			Path: "/api/v2/tokens/",
		},
		PageSize: pointer.To(100),
	})
	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return tokens, nil
}

// Create creates a new token.
func (t *TokenService) Create() (*Token, error) {
	if t.backref.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}
	token, err := doPost[Token](t.backref, &PostRequest{
		Request: Request{
			Path: "/api/v2/tokens/",
		},
	})
	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return token, nil
}

// Delete deletes one or more tokens.
func (t *TokenService) Delete(tokens ...string) (*Token, error) {
	if t.backref.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}
	token, err := doPost[Token](t.backref, &PostRequest{
		Request: Request{
			Path: "/api/v2/tokens/",
		},
	})
	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return token, nil
}
