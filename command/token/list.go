package token

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
)

type List struct {
	base.Command
	// Username string `short:"u" long:"username" description:"The username to use for authentication." required:"yes" env:"SMS_USERNAME" cfg:"username"`
	// Password string `short:"p" long:"password" description:"The password to use for authentication." required:"yes" env:"SMS_PASSWORD" cfg:"password"`
}

// Execute is the real implementation of the token.List command.
func (cmd *List) Execute(args []string) error {
	slog.Debug("called token list command")

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

	tokens, err := client.Token.List()
	if err != nil {
		slog.Error("error performing token list API call", "error", err)
		fmt.Printf("error: %s\n", color.RedString(err.Error()))
		return fmt.Errorf("error performing API call: %w", err)
	}

	for _, token := range tokens {
		expiry := token.ExpiryDate
		if expiry.IsZero() {
			fmt.Printf("token: %s (no expiration)\n", color.YellowString(token.Token))
		} else {
			fmt.Printf("token: %s (expires on %s)\n", color.YellowString(token.Token), color.YellowString(expiry.Format(time.RFC3339)))
		}
	}
	return nil
}
