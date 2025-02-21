package rdcom

import (
	"crypto/tls"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"resty.dev/v3"
)

// AuthType represents the available authentication types.
type AuthType int8

// List of available authentication types.
const (
	BasicAuth AuthType = iota
	TokenAuth
)

// Client is the API client.
type Client struct {
	// api is the underlying API client using Resty.
	api *resty.Client
	// agent is the user agent to use in HTTP requests to the API.
	agent string
	// endpoint is the base URL for API requests.
	endpoint string `validate:"required"`
	// username is used with password to perform basic authentication.
	username string
	// password is used to perform basic authentication.
	password string
	// token is the authentication token to use in requests.
	token string
	// account is the ID of the account to use ito scope requests.
	account string `validate:"required"`
	// TokenService is the Token service
	TokenService *TokenService
}

// Service represents an API service.
type Service struct {
	client *Client
}

// option allows to set options in a functional way.
type Option func(*Client)

// WithUserAgent sets the API client user agent
func WithUserAgent(agent string) Option {
	return func(c *Client) {
		slog.Debug("setting user agent", "agent", agent)
		c.agent = agent
		c.api.SetHeader("User-Agent", agent)
	}
}

// WithBaseURL sets the API endpoint.
func WithBaseURL(endpoint string) Option {
	return func(c *Client) {
		slog.Debug("setting base URL", "endpoint", endpoint)
		c.endpoint = endpoint
		c.api.SetBaseURL(endpoint)
	}
}

// WithAccount sets the account option.
func WithAccount(account string) Option {
	return func(c *Client) {
		slog.Debug("setting account", "account", account)
		c.account = account
	}
}

// WithSkipTLSVerify sets the skip TLS verification option.
func WithSkipTLSVerify(insecure bool) Option {
	return func(c *Client) {
		slog.Debug("setting skip TLS verification", "insecure", insecure)
		c.api.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: insecure,
		})
		c.api.SetDisableWarn(true)
	}
}

// WithDebug sets the debug option.
func WithDebug() Option {
	return func(c *Client) {
		slog.Debug("enabling debug")
		c.api.EnableDebug()
	}
}

// WithTrace sets the trace option.
func WithTrace() Option {
	return func(c *Client) {
		slog.Debug("enabling trace")
		c.api.EnableTrace()
	}
}

// WithUserCredentials sets the basic authentication credentials for the user.
func WithUserCredentials(username string, password string) Option {
	return func(c *Client) {
		slog.Debug("setting user credentials", "username", username, "password", password)
		c.username = username
		c.password = password
		c.api.SetBasicAuth(username, password)
	}
}

// WithAuthToken sets the authentication token for the user.
func WithAuthToken(token string) Option {
	return func(c *Client) {
		slog.Debug("setting authentication token", "token", token)
		c.token = token
		c.api.SetAuthToken(token)
	}
}

// New creates anew API client.
func New(options ...Option) *Client {

	c := &Client{
		api: resty.New(),
	}
	for _, option := range options {
		option(c)
	}
	c.TokenService = &TokenService{Service{client: c}}
	// TODO: initialise more services here...

	// perform struct level validation
	validate := validator.New()
	validate.RegisterStructValidation(check, Client{})
	if err := validate.Struct(c); err != nil {
		slog.Error("invalid API client configuration", "error", err)
		return nil
	}
	slog.Debug("API client ready")
	return c
}

// Close frees the API client resources.
func (c *Client) Close() error {
	return c.api.Close()
}

// Options is a struct that contains the outline of an API call
// options.
type Options struct {
	EntityPath  string            `json:"entity_path" yaml:"entity_path" validate:"required"`
	PathParams  map[string]string `json:"path_params" validate:"required"`
	QueryParams map[string]string `json:"query_params" validate:"required"`
}

type GetOptions Options

// Get performs an API request to retrieve one single entity.
func Get[T any](client *Client, options *GetOptions) (*T, error) {
	var err error
	request := client.api.R()

	if options.QueryParams != nil {
		slog.Debug("setting query params", "values", options.QueryParams)
		request.SetQueryParams(options.QueryParams)
	}

	if options.PathParams != nil {
		slog.Debug("setting path params", "values", options.PathParams)
		request.SetPathParams(options.PathParams)
	}

	// if options.Input != nil {
	// 	slog.Debug("setting input struct", "type", fmt.Sprintf("%T", options.Output))
	// 	request.SetBody(options.Input)
	// } else {
	// 	slog.Debug("no input struct provided")
	// }

	// if options.Output != nil {
	// 	slog.Debug("setting output struct", "type", fmt.Sprintf("%T", options.Output))
	// 	request.SetResult(options.Output)
	// } else {
	// 	slog.Warn("no output struct provided")
	// }

	result := new(T)
	request.SetResult(result)
	_, err = request.Get(options.EntityPath)

	// switch options.Method {
	// case http.MethodGet:
	// 	_, err = request.Get(call.Path)
	// case http.MethodPost:
	// 	_, err = request.Post(call.Path)
	// case http.MethodDelete:
	// 	_, err = request.Delete(call.Path)
	// default:
	// 	slog.Error("unsupported method", "method", call.Method)
	// 	return nil, fmt.Errorf("unsupported method: %s", call.Method)
	// }

	if err != nil {
		slog.Error("error performing GET API request", "entity", options.EntityPath, "error", err)
		return nil, err
	}

	slog.Debug("GET API options successful", "path", options.EntityPath)
	return result, nil
}

type ListOptions struct {
	Options  `json:",inline"`
	PageSize *int `json:"page_size,omitempty" yaml:"page_size,omitempty"`
}

// List performs an API request to retrieve multiple entities, possibly using pagination.
func List[T any](client *Client, options *ListOptions) ([]T, error) {
	request := client.api.R()

	if options.QueryParams != nil {
		slog.Debug("setting query params", "values", options.QueryParams)
		request.SetQueryParams(options.QueryParams)
	}

	if options.PathParams != nil {
		slog.Debug("setting path params", "values", options.PathParams)
		request.SetPathParams(options.PathParams)
	}

	type payload[T any] struct {
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

	results := make([]T, 0)
	page := &payload[T]{}
	offset := 0
	for {
		if options.PageSize != nil {
			slog.Debug("enabling pagination", "page size", *options.PageSize)
			request.SetQueryParam("paginated-view", "true")
			request.SetQueryParam("limit", fmt.Sprintf("%d", *options.PageSize))
			request.SetQueryParam("offset", fmt.Sprintf("%d", offset))
			offset = offset + 1
		}
		_, err := request.
			SetResult(page).
			Get(options.EntityPath)

		if err != nil {
			slog.Error("error performing GET (many) API request", "error", err)
			return nil, err
		}
		slog.Debug("API call successful", "count", len(page.Results))
		results = append(results, page.Results...)

		if page.TotPages == 1 || offset >= page.TotPages {
			slog.Debug("no more pages")
			break
		}
	}
	return results, nil
}

type CreateOptions Options

// Create performs an API request to create a new entity.
func Create[T any](client *Client, entity *T, options *CreateOptions) (*T, error) {
	request := client.api.R()

	if entity != nil {
		slog.Debug("setting entity", "type", fmt.Sprintf("%T", entity), "value", *entity)
		request.SetBody(entity)
	} else {
		slog.Warn("no entity provided?")
	}

	result := new(T)
	_, err := request.
		SetResult(result).
		Post(options.EntityPath)

	if err != nil {
		slog.Error("error performing POST API request", "error", err)
		return nil, err
	}
	slog.Debug("API call success", "result", result)
	return result, nil
}

func check(sl validator.StructLevel) {

	client := sl.Current().Interface().(Client)

	if client.token != "" && (client.username != "" || client.password != "") {
		slog.Error("cannot have both token and username/password")
		sl.ReportError(client.token, "token", "Token", "tokenorcreds", "")
	}

	if client.token == "" {
		if client.username == "" || client.password == "" {
			slog.Error("must have both username and password")
			sl.ReportError(client.token, "username", "Username", "validcreds", "")
			sl.ReportError(client.token, "password", "Password", "validcreds", "")
		}
	}

	if client.endpoint == "" {
		slog.Error("must have an endpoint")
		sl.ReportError(client.endpoint, "endpoint", "Endpoint", "required", "")
	}

	if client.TokenService.client == nil {
		slog.Error("token service not initialised")
		sl.ReportError(client.TokenService.client, "client", "Token.Service", "service", "")
	}

	// add more services here...
}
