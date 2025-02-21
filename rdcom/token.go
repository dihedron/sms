package rdcom

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
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
	if t.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	call := &Call[any, ListResponse[Token]]{
		Method: http.MethodGet,
		Path:   "/api/v2/tokens",
		Output: &ListResponse[Token]{},
	}

	result, err := Do(t.client, call)

	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return result.Results, nil
}

// Create creates a new token.
func (t *TokenService) Create() (*Token, error) {
	if t.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}
	token, err := doPost[Token](t.client, &PostRequest{
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
func (t *TokenService) Delete(token string) (*Token, error) {

	if t.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	response, err := doDelete[Token](t.client, &DeleteRequest[Token]{
		Request: Request{
			Path: "/api/v2/tokens/",
		},
		Value: Token{
			Token: token,
		},
	})
	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return response, nil
}
