package ping

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
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
		slog.Error("error performing token list API call", "error", err)
		fmt.Printf("connection: %s\n", color.RedString("KO"))
		return fmt.Errorf("error performing API call: %w", err)
	}

	slog.Debug("successful token list API call")
	fmt.Printf("connection: %s\n", color.GreenString("OK"))
	return nil
}
