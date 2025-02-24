package ping

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
)

type Ping struct {
	base.TokenCommand
	// Account is the account to use in APIs invocation.
	// Account string `short:"a" long:"account" description:"The account to use for authentication." required:"yes" env:"SMS_ACCOUNT" cfg:"account"`
}

// Execute is the real implementation of the Ping command.
func (cmd *Ping) Execute(args []string) error {
	slog.Debug("called ping command", "token", cmd.Token, "endpoint", cmd.Endpoint)

	options := []rdcom.Option{
		rdcom.WithBaseURL(cmd.Endpoint),
		rdcom.WithUserAgent("bancaditalia/0.1"),
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
	if cmd.Token != nil {
		options = append(options, rdcom.WithAuthToken(*cmd.Token))
	}

	client, err := rdcom.New(options...)
	if err != nil {
		slog.Error("error initialising API client", "error", err)
		return err
	}

	defer client.Close()

	if _, err := client.TokenService.List(); err != nil {
		slog.Error("error performing token list API call", "error", err)
		fmt.Printf("connection: %s\n", color.RedString("KO"))
		return fmt.Errorf("error performing API call: %w", err)
	}

	slog.Debug("successful token list API call")
	fmt.Printf("connection: %s\n", color.GreenString("OK"))
	return nil
}
