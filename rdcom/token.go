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
	TotPages               int     `json:"tot_pages"`
	CurrentPageFirstRecord int     `json:"current_page_first_record"`
	CurrentPageLastRecord  int     `json:"current_page_last_record"`
	Limit                  int     `json:"limit"`
	Offset                 int     `json:"offset"`
	Count                  int     `json:"count"`
	CountIsEstimate        bool    `json:"count_is_estimate"`
	Next                   *string `json:"next"`
	Previous               *string `json:"previous"`
	Results                []Token `json:"results"`
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

	tokens, err := doQuery(t.backref, &Query{
		Path:     "/api/v2/tokens/",
		PageSize: pointer.To(1),
	})
	if err != nil {
		return nil, err
	}

	// tokens := []Token{}

	// result := &TokenListResponse{}
	// response, err := t.backref.api.
	// 	R().
	// 	SetQueryParam("paginated-view", fmt.Sprintf("%d", page)).
	// 	SetResult(result).
	// 	Get("/api/v2/tokens/")

	// if err != nil {
	// 	slog.Error("error performing GET API request", "error", err)
	// 	return nil, err
	// }

	// slog.Debug("API response", "response", response)

	// tokens = append(tokens, result.Results...)

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
