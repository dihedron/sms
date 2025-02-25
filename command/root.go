package command

import (
	"github.com/dihedron/sms/command/account"
	"github.com/dihedron/sms/command/ping"
	smsgateway "github.com/dihedron/sms/command/sms_gateway"
	"github.com/dihedron/sms/command/token"
	"github.com/dihedron/sms/command/version"
)

// Commands is the set of root command groups.
type Commands struct {

	// Check checks the connectivity to RDCom API.
	Ping ping.Ping `command:"ping" alias:"p" description:"Try to connect to the RDCom API server."`

	// Account is a subcommand group related to account management.
	//lint:ignore SA5008 commands can have multiple aliases
	Account account.Account `command:"account" alias:"acc" alias:"a" description:"Account-related operations."`

	// SMSGateway is a subcommand group related to SMS gateway management.
	//lint:ignore SA5008 commands can have multiple aliases
	SMSGateway smsgateway.SMSGateway `command:"sms_gateway" alias:"smsgw" alias:"gw" alias:"g" description:"SMS gateway-related operations."`

	// Token is a subcommand group related to token management.
	//lint:ignore SA5008 commands can have multiple aliases
	Token token.Token `command:"token" alias:"tok" alias:"tk" alias:"t" description:"Token management operations."`

	// Version prints the application version information and exits.
	//lint:ignore SA5008 commands can have multiple aliases
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit."`
}
