package command

import (
	"github.com/dihedron/sms/command/check"
	"github.com/dihedron/sms/command/version"
)

// Commands is the set of root command groups.
type Commands struct {

	// Check checks the connectivity to RDCom API.
	Check check.Check `command:"check" alias:"c" description:"Try to connect to the RDCom API server."`

	// Version prints the application version information and exits.
	//lint:ignore SA5008 commands can have multiple aliases
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit."`
}
