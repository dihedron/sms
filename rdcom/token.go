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
	Token      string    `json:"token,omitzero"`
	ExpiryDate time.Time `json:"expire_date,omitzero"`
}

// List returns the list of tokens.
func (t *TokenService) List() ([]Token, error) {
	if t.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	options := &ListOptions{
		Options: Options{
			EntityPath: "/api/v2/tokens",
		},
		PageSize: pointer.To(100),
	}

	result, err := List[Token](t.client, options)

	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return result, nil
}

// Create creates a new token.
func (t *TokenService) Create() (*Token, error) {
	if t.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}
	token, err := Create[Token](t.client, nil, &CreateOptions{
		EntityPath: "/api/v2/tokens/",
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
