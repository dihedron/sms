package token

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
)

type Create struct {
	//base.CredentialsCommand
	base.TokenCommand
}

// Execute is the real implementation of the token.List command.
func (cmd *Create) Execute(args []string) error {
	slog.Debug("called token create command")

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

	client := rdcom.New(options...)

	defer client.Close()

	token, err := client.TokenService.Create()
	if err != nil {
		slog.Error("error performing token create API call", "error", err)
		fmt.Printf("error: %s\n", color.RedString(err.Error()))
		return fmt.Errorf("error performing API call: %w", err)
	}

	expiry := token.ExpiryDate
	if expiry.IsZero() {
		fmt.Printf("token: %s (no expiration)\n", color.YellowString(token.Token))
	} else {
		fmt.Printf("token: %s (expires on %s)\n", color.YellowString(token.Token), color.YellowString(expiry.Format(time.RFC3339)))
	}
	return nil
}
