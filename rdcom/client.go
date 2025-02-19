package rdcom

import (
	"crypto/tls"

	"resty.dev/v3"
)

// Client is the API client.
type Client struct {
	api      *resty.Client
	username string
	password string
	Token    *Token
}

// Service represents an API service.
type Service struct {
	backref *Client
}

// New creates anew API client.
func New(endpoint string, insecure bool) *Client {
	c := &Client{
		api: resty.
			New().
			SetBaseURL(endpoint).
			SetTLSClientConfig(&tls.Config{
				InsecureSkipVerify: insecure,
			}).
			EnableTrace().
			EnableDebug(),
		Token: &Token{},
		// more services here
	}
	c.Token.backref = c
	// add backrefs here
	return c
}

// Close frees the API client resources.
func (c *Client) Close() error {
	return c.api.Close()
}

// SetUserCredentials sets the basic authentication credentials for the user.
func (c *Client) SetUserCredentials(username string, password string) *Client {
	c.username = username
	c.password = password
	return c
}
