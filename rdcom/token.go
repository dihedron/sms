package rdcom

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/sms/pointer"
)

type TokenService struct {
	Service
}

func (t *TokenService) Create() error {
	if t.backref.username == "" || t.backref.password == "" {
		slog.Error("invalid credentials")
		return errors.New("invalid credentials")
	}

	response, err := t.Service.backref.api.
		R().
		SetBasicAuth(t.backref.username, t.backref.username).
		Post("/api/v2/tokens/")
	if err != nil {
		slog.Error("error performing POST API request", "error", err)
		return err
	}

	slog.Debug("API response", "response", response)
	return nil
}

type TokenListResponse struct {
	ListResponse `json:",inline"`
	Tokens       []Token `json:"results"`
}

type Token struct {
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expire_date"`
}

func (t *TokenService) List() ([]Token, error) {

	if t.backref.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	tokens, err := doQuery(t.backref, &ListRequest{
		Path:     "/api/v2/tokens/",
		PageSize: pointer.To(1),
	})
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		expiry := token.ExpiryDate
		if expiry.IsZero() {
			fmt.Printf("token: %s (no expiration)\n", token.Token)
		} else {
			fmt.Printf("token: %s (expires on %s)\n", token.Token, expiry.Format(time.RFC3339))
		}
	}

	return tokens, nil
}
