package version

import (
	"log/slog"
	"os"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/version"
)

// Version is the command that prints information about the application
// or plugin to the console; it support both compact and verbose mode.
type Version struct {
	base.Command
}

// Execute is the real implementation of the Version command.
func (cmd *Version) Execute(args []string) error {
	slog.Debug("running version command")
	version.Print(os.Stdout)
	slog.Debug("command done")
	return nil
}
