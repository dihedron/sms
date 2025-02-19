package rdcom

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
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

type ListResponse struct {
	TotPages               int  `json:"tot_pages"`
	CurrentPageFirstRecord int  `json:"current_page_first_record"`
	CurrentPageLastRecord  int  `json:"current_page_last_record"`
	Limit                  int  `json:"limit"`
	Offset                 int  `json:"offset"`
	Count                  int  `json:"count"`
	CountIsEstimate        bool `json:"count_is_estimate"`
	Next                   int  `json:"next"`
	Previous               int  `json:"previous"`
	Results                []struct {
		Token      string    `json:"token"`
		ExpireDate time.Time `json:"expire_date"`
	} `json:"results"`
}

func (t *Token) List() ([]string, error) {

	if t.backref.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	result := &ListResponse{}

	response, err := t.Service.backref.api.
		R().
		SetResult(result).
		SetAuthToken(t.backref.token).
		Get("/api/v2/tokens/")
	if err != nil {
		slog.Error("error performing GET API request", "error", err)
		return nil, err
	}

	for _, token := range result.Results {
		expiry := token.ExpireDate
		if expiry.IsZero() {
			fmt.Printf("token: %s (no expiration)\n", token.Token)
		} else {
			fmt.Printf("token: %s (expires on %s)\n", token.Token, expiry.Format(time.RFC3339))
		}
	}

	slog.Debug("API response", "response", response)

	return nil, nil
}
