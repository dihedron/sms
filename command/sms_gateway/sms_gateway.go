package smsgateway

type SMSGateway struct {
	// // Create is the command to create a new token.
	// //lint:ignore SA5008 commands can have multiple aliases
	// Create Create `command:"create" alias:"cr" alias:"c" description:"Create a new token."`

	// // Delete is the command to delete a token.
	// //lint:ignore SA5008 commands can have multiple aliases
	// Delete Delete `command:"delete" alias:"del" alias:"d" description:"Delete a token."`

	// List is the command to list accounts.
	//lint:ignore SA5008 commands can have multiple aliases
	List List `command:"list" alias:"ls" alias:"l" description:"List existing accounts."`
}
