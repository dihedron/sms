package check

import (
	"log/slog"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/rdcom"
)

type Check struct {
	base.Command

	Username string `short:"u" long:"username" description:"The username to use for authentication." required:"yes" env:"SMS_USERNAME" cfg:"username"`
	Password string `short:"p" long:"password" description:"The password to use for authentication." required:"yes" env:"SMS_PASSWORD" cfg:"password"`
	Endpoint string `short:"e" long:"endpoint" description:"The API endpoint to use." required:"yes" env:"SMS_ENDPOINT" cfg:"endpoint" default:"https://platform.rdcom.com"`
}

// Execute is the real implementation of the Version command.
func (cmd *Check) Execute(args []string) error {
	slog.Debug("called check command")

	client := rdcom.New(cmd.Endpoint, true).SetUserCredentials(cmd.Username, cmd.Password)
	defer client.Close()

	if err := client.Token.Create(); err != nil {
		slog.Error("error creating login token", "error", err)
		return err
	}

	return nil
}
