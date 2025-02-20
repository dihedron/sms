package rdcom

import (
	"fmt"
	"log/slog"
)

type Request struct {
	Path string `json:"path" validate:"required"`
}

type GetRequest struct {
	Request     `json:",inline"`
	PathParams  map[string]string `json:"path_params" validate:"required"`
	QueryParams map[string]string `json:"query_params" validate:"required"`
	PageSize    *int              `json:"page_size,omitempty"`
}

type GetResponse[T any] struct {
	TotPages               int     `json:"tot_pages"`
	CurrentPageFirstRecord int     `json:"current_page_first_record"`
	CurrentPageLastRecord  int     `json:"current_page_last_record"`
	Limit                  int     `json:"limit"`
	Offset                 int     `json:"offset"`
	Count                  int     `json:"count"`
	CountIsEstimate        bool    `json:"count_is_estimate"`
	Next                   *string `json:"next"`
	Previous               *string `json:"previous"`
	Results                []T     `json:"results"`
}

func doGet[T any](client *Client, info *GetRequest) ([]T, error) {
	results := make([]T, 0)
	request := client.api.R()
	if info.QueryParams != nil {
		slog.Debug("setting query params", "values", info.QueryParams)
		request.SetQueryParams(info.QueryParams)
	}
	if info.PathParams != nil {
		slog.Debug("setting path params", "values", info.PathParams)
		request.SetPathParams(info.PathParams)
	}
	result := &GetResponse[T]{}
	offset := 0
	for {
		if info.PageSize != nil {
			slog.Debug("enabling pagination", "page size", info.PageSize)
			request.SetQueryParam("paginated-view", "true")
			request.SetQueryParam("limit", fmt.Sprintf("%d", *info.PageSize))
			request.SetQueryParam("offset", fmt.Sprintf("%d", offset))
			offset = offset + 1
		}
		_, err := request.
			SetResult(result).
			Get(info.Path)

		if err != nil {
			slog.Error("error performing GET API request", "error", err)
			return nil, err
		}
		results = append(results, result.Results...)

		if result.TotPages == 1 || offset >= result.TotPages {
			slog.Debug("no more pages")
			break
		}
	}
	return results, nil
}

type PostRequest struct {
	Request `json:",inline"`
}

func doPost[T any](client *Client, info *PostRequest) (*T, error) {
	request := client.api.R()
	result := new(T)

	_, err := request.
		SetResult(result).
		Post(info.Path)

	if err != nil {
		slog.Error("error performing POST API request", "error", err)
		return nil, err
	}
	return result, nil
}
