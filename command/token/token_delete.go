package token

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
)

type Delete struct {
	base.TokenCommand
}

// Execute is the real implementation of the token delete command.
func (cmd *Delete) Execute(args []string) error {
	slog.Debug("called token delete command", "args", args)

	if len(args) == 0 {
		slog.Error("no token ID provided")
		return fmt.Errorf("no token ID provided")
	}

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

	// if cmd.Username != "" && cmd.Password != "" {
	// 	options = append(options, rdcom.WithUserCredentials(cmd.Username, cmd.Password))
	// }

	client, err := rdcom.New(options...)
	if err != nil {
		slog.Error("error initialising API client", "error", err)
		return err
	}

	defer client.Close()

	for _, arg := range args {

		token, err := client.TokenService.Delete(arg)
		if err != nil {
			slog.Error("error performing token delete API call", "error", err)
			return err
		}

		expiry := token.ExpiryDate
		if expiry.IsZero() {
			fmt.Printf("token: %s (no expiration)\n", color.YellowString(token.Token))
		} else {
			fmt.Printf("token: %s (expires on %s)\n", color.YellowString(token.Token), color.YellowString(expiry.Format(time.RFC3339)))
		}
	}
	return nil
}
