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
	result := &ListResponse[T]{}
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

type DeleteRequest[T any] struct {
	Request `json:",inline"`
	Value   T `json:"value"`
}

func doDelete[T any](client *Client, info *DeleteRequest[T]) (*T, error) {
	request := client.api.R()
	result := new(T)
	_, err := request.
		SetResult(result).
		SetBody(info.Value).
		Delete(info.Path)

	if err != nil {
		slog.Error("error performing POST API request", "error", err)
		return nil, err
	}
	return result, nil
}

/*
type Call[I any, O any] struct {
	Method      string            `json:"method" yaml:"method" validate:"required"`
	Path        string            `json:"path" yaml:"path" validate:"required"`
	PathParams  map[string]string `json:"path_params" validate:"required"`
	QueryParams map[string]string `json:"query_params" validate:"required"`
	PageSize    *int              `json:"page_size,omitempty" yaml:"page_size,omitempty"`
	Input       *I                `json:"input" yaml:"input"`
	Output      *O                `json:"output" yaml:"output"`
}


func doRun[S any, T any](client *Client, call *Call[S, T]) error {

	request := client.api.R()
	if call.QueryParams != nil {
		slog.Debug("setting query params", "values", call.QueryParams)
		request.SetQueryParams(call.QueryParams)
	}
	if call.PathParams != nil {
		slog.Debug("setting path params", "values", call.PathParams)
		request.SetPathParams(call.PathParams)
	}

	offset := 0
	for {
		if call.Method == http.MethodGet && call.PageSize != nil {
			slog.Debug("enabling pagination", "page size", call.PageSize)
			request.SetQueryParam("paginated-view", "true")
			request.SetQueryParam("limit", fmt.Sprintf("%d", *call.PageSize))
			request.SetQueryParam("offset", fmt.Sprintf("%d", offset))
			offset = offset + 1
		}

		if
		_, err := request.SetResult(call.Output).
			Get(call.Path)

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

	return nil
}
*/
