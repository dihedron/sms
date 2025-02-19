package command

import (
	"github.com/dihedron/sms/command/ping"
	"github.com/dihedron/sms/command/version"
)

// Commands is the set of root command groups.
type Commands struct {

	// Check checks the connectivity to RDCom API.
	Ping ping.Ping `command:"ping" alias:"p" description:"Try to connect to the RDCom API server."`

	Account struct {
		Show string `command:"show" alias:"s" description:"Show account details."`

		// Create is the command to create a new account.
		Create string `command:"create" alias:"c" description:"Create a new account."`

		// Delete is the command to delete an account.
		Delete string `command:"delete" alias:"d" description:"Delete an account."`

		// List is the command to list accounts.
		List string `command:"list" alias:"l" description:"List accounts."`

		// Refresh is the command to refresh an account.
		Refresh string `command:"refresh" alias:"r" description:"Refresh an account."`
	} `command:"account" alias:"a" description:"Account-related operations."`

	Token struct {
		// Create is the command to create a new token.
		Create string `command:"create" alias:"c" description:"Create a new token."`

		// Delete is the command to delete a token.
		Delete string `command:"delete" alias:"d" description:"Delete a token."`

		// List is the command to list tokens.
		List string `command:"list" alias:"l" description:"List tokens."`

		// Refresh is the command to refresh a token.
		Refresh string `command:"refresh" alias:"r" description:"Refresh a token."`
	} `command:"token" alias:"t" description:"Token-related operations."`

	// Version prints the application version information and exits.
	//lint:ignore SA5008 commands can have multiple aliases
	Version version.Version `command:"version" alias:"ver" alias:"v" description:"Show the command version and exit."`
}
