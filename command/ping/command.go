package ping

import (
	"log/slog"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
)

type Ping struct {
	base.Command

	Account string `short:"a" long:"account" description:"The account to use for authentication." required:"yes" env:"SMS_ACCOUNT" cfg:"account"`

	Token string `short:"t" long:"token" description:"The token to use for authentication." required:"yes" env:"SMS_TOKEN" cfg:"token"`

	Endpoint string `short:"e" long:"endpoint" description:"The API endpoint to use." required:"yes" env:"SMS_ENDPOINT" cfg:"endpoint" default:"https://platform.rdcom.com"`
}

// Execute is the real implementation of the Ping command.
func (cmd *Ping) Execute(args []string) error {
	slog.Debug("called ping command", "token", cmd.Token, "endpoint", cmd.Endpoint, "account", cmd.Account)

	options := []rdcom.Option{
		rdcom.WithBaseURL(cmd.Endpoint),
	}
	if cmd.SkipVerifyTLS {
		options = append(options, rdcom.WithSkipTLSVerify(true))
	}
	if cmd.EnableDebug {
		options = append(options, rdcom.WithDebug())
	}
	if cmd.EnableTrace {
		options = append(options, rdcom.WithTrace())
	}
	if cmd.Token != "" {
		options = append(options, rdcom.WithAuthToken(cmd.Token))
	}

	client := rdcom.New(options...)

	defer client.Close()

	if _, err := client.Token.List(); err != nil {
		slog.Error("error listing tokens", "error", err)
		return err
	}

	slog.Debug("tokens listed")
	return nil
}
