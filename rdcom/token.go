package rdcom

import (
	"errors"
	"log/slog"
)

type Token struct {
	Service
}

func (t *Token) Create() error {
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

func (t *Token) List() ([]string, error) {

	if t.backref.username == "" || t.backref.password == "" {
		slog.Error("invalid credentials")
		return nil, errors.New("invalid credentials")
	}

	response, err := t.Service.backref.api.
		R().
		SetBasicAuth(t.Service.backref.username, t.Service.backref.username).
		Get("/api/v2/tokens/")
	if err != nil {
		slog.Error("error performing GET API request", "error", err)
		return nil, err
	}

	slog.Debug("API response", "response", response)

	return nil, nil
}
