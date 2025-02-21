package rdcom

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"

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

// Call is a struct that represents a specific API call, along with its optional
// input data (to be placed in the request body) and its result.
type Call[I any, O any] struct {
	Method      string            `json:"method" yaml:"method" validate:"required"`
	Path        string            `json:"path" yaml:"path" validate:"required"`
	PathParams  map[string]string `json:"path_params" validate:"required"`
	QueryParams map[string]string `json:"query_params" validate:"required"`
	PageSize    *int              `json:"page_size,omitempty" yaml:"page_size,omitempty"`
	Input       *I                `json:"input" yaml:"input"`
	Output      *O                `json:"output" yaml:"output"`
}

// ListResponse is the response to List requests.
type ListResponse[T any] struct {
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

// Do performs an API request.
func Do[I any, O any](client *Client, call *Call[I, O]) (*O, error) {
	var err error
	request := client.api.R()

	if call.QueryParams != nil {
		slog.Debug("setting query params", "values", call.QueryParams)
		request.SetQueryParams(call.QueryParams)
	}

	if call.PathParams != nil {
		slog.Debug("setting path params", "values", call.PathParams)
		request.SetPathParams(call.PathParams)
	}

	if call.Input != nil {
		slog.Debug("setting input struct", "type", fmt.Sprintf("%T", call.Output))
		request.SetBody(call.Input)
	} else {
		slog.Debug("no input struct provided")
	}

	if call.Output != nil {
		slog.Debug("setting output struct", "type", fmt.Sprintf("%T", call.Output))
		request.SetResult(call.Output)
	} else {
		slog.Warn("no output struct provided")
	}

	switch call.Method {
	case http.MethodGet:
		_, err = request.Get(call.Path)
	case http.MethodPost:
		_, err = request.Post(call.Path)
	case http.MethodDelete:
		_, err = request.Delete(call.Path)
	default:
		slog.Error("unsupported method", "method", call.Method)
		return nil, fmt.Errorf("unsupported method: %s", call.Method)
	}

	if err != nil {
		slog.Error("error performing API request", "method", call.Method, "path", call.Path, "error", err)
		return nil, err
	}

	slog.Debug("API call successful", "method", call.Method, "path", call.Path)
	return call.Output, nil
}

// Close frees the API client resources.
func (c *Client) Close() error {
	return c.api.Close()
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
