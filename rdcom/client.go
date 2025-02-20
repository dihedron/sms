package rdcom

import (
	"crypto/tls"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"resty.dev/v3"
)

// Client is the API client.
type Client struct {
	api      *resty.Client
	endpoint string `validate:"required"`
	username string
	password string
	account  string `validate:"required"`
	token    string
	Token    *TokenService
}

// Service represents an API service.
type Service struct {
	backref *Client
}

type Option func(*Client)

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
	c.Token = &TokenService{Service{backref: c}}
	// initialise more services here...

	// perform struct level validation
	validate := validator.New()
	validate.RegisterStructValidation(check, Client{})
	if err := validate.Struct(c); err != nil {
		slog.Error("invalid API client configuration", "error", err)
		return nil
	}

	return c
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

	if client.Token.backref == nil {
		slog.Error("token service not initialised")
		sl.ReportError(client.Token.backref, "backref", "Token.Service", "service", "")
	}

	// add more services here...
}
