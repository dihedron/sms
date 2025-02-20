package rdcom

import (
	"fmt"
	"log/slog"
	"time"
)

type Query struct {
	Path        string            `json:"path" validate:"required"`
	PathParams  map[string]string `json:"path_params" validate:"required"`
	QueryParams map[string]string `json:"query_params" validate:"required"`
	PageSize    *int              `json:"page_size,omitempty"`
}

func doQuery(client *Client, query *Query) ([]Token, error) {

	tokens := []Token{}

	request := client.api.R()

	if query.QueryParams != nil {
		slog.Debug("setting query params", "values", query.QueryParams)
		request.SetQueryParams(query.QueryParams)
	}

	if query.PathParams != nil {
		slog.Debug("setting path params", "values", query.PathParams)
		request.SetPathParams(query.PathParams)
	}

	result := &TokenListResponse{}

	offset := 0
	for {
		if query.PageSize != nil {
			slog.Debug("enabling pagination", "page size", query.PageSize)
			request.SetQueryParam("paginated-view", "true")
			request.SetQueryParam("limit", fmt.Sprintf("%d", *query.PageSize))
			request.SetQueryParam("offset", fmt.Sprintf("%d", offset))
			offset = offset + 1
		}

		_, err := request.
			SetResult(result).
			Get("/api/v2/tokens/")

		if err != nil {
			slog.Error("error performing GET API request", "error", err)
			return nil, err
		}
		tokens = append(tokens, result.Results...)

		if result.TotPages == 1 || offset >= result.TotPages {
			slog.Debug("no more pages")
			break
		}
		time.Sleep(1 * time.Second)
	}
	return tokens, nil
}
