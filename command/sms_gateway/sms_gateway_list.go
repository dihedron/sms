package smsgateway

import (
	"fmt"
	"log/slog"

	"github.com/dihedron/sms/command/base"
	"github.com/dihedron/sms/format"
	"github.com/dihedron/sms/rdcom"
	"github.com/fatih/color"
)

// List is the SMS gateway list command.
type List struct {
	base.TokenCommand
	// Account is the account whose SMS gateways to list.
	Account string `short:"a" long:"account" description:"The account whose SMS gateways to list." required:"yes" env:"SMS_ACCOUNT"`
}

// Execute is the real implementation of the account list command.
func (cmd *List) Execute(args []string) error {
	slog.Debug("called SMS gateway list command")

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

	gateways, err := client.SMSGatewayService.List(cmd.Account)
	if err != nil {
		slog.Error("error performing token list API call", "error", err)
		fmt.Printf("error: %s\n", color.RedString(err.Error()))
		return fmt.Errorf("error performing API call: %w", err)
	}

	slog.Info("gateways", "length", len(gateways))

	for _, gateway := range gateways {
		fmt.Printf("gateway: %s\n", color.YellowString(fmt.Sprintf("%d", gateway.ID)))
		fmt.Printf(" - type                   : %s\n", color.YellowString(gateway.GatewayType))
		// fmt.Printf(" - parent                 : %s\n", color.YellowString(gateway.Parent))
		// fmt.Printf(" - enabled                : %s\n", format.ColoredBool(gateway.Enabled))
		// fmt.Printf(" - creation               : %s\n", color.YellowString(gateway.Created.Format(base.DefaultDateFormat)))
		// if account.ExpirationDate.IsZero() {
		// 	fmt.Printf(" - expiration             : %s\n", color.YellowString("none"))
		// } else {
		// 	fmt.Printf(" - expiration             : %s\n", color.YellowString(account.ExpirationDate.Format(base.DefaultDateFormat)))
		// }
		// fmt.Printf(" - main contact           :\n")
		// fmt.Printf("   - name                 : %s\n", color.YellowString(account.Infos.MainContactName))
		// fmt.Printf("   - name                 : %s\n", color.YellowString(account.Infos.MainContactSurname))
		// fmt.Printf("   - email                : %s\n", color.YellowString(account.Infos.MainContactEmail))
		// fmt.Printf("   - mobile               : %s\n", color.YellowString(account.Infos.MainContactCell))
		// fmt.Printf(" - representative         :\n")
		// fmt.Printf("   - name                 : %s\n", color.YellowString(account.Infos.ReprName))
		// fmt.Printf("   - name                 : %s\n", color.YellowString(account.Infos.ReprSurname))
		// fmt.Printf("   - email                : %s\n", color.YellowString(account.Infos.ReprEmail))
		// fmt.Printf(" - company                :\n")
		// fmt.Printf("   - name                 : %s\n", color.YellowString(account.Infos.Company))
		// fmt.Printf("   - address              : %s\n", color.YellowString(account.Infos.Address))
		// fmt.Printf("   - city                 : %s\n", color.YellowString(account.Infos.City))
		// fmt.Printf("   - state                : %s\n", color.YellowString(account.Infos.State))
		// fmt.Printf("   - country              : %s\n", color.YellowString(account.Infos.Country))
		// fmt.Printf("   - ZIP code             : %s\n", color.YellowString(account.Infos.ZipCode))
		// fmt.Printf(" - enablements            :\n")
		// fmt.Printf("   - email preview        : %s\n", format.ColoredBool(account.EnableEmailPreview))
		// fmt.Printf("   - PDF attachments      : %s\n", format.ColoredBool(account.EnablePdfAttachments))
		// fmt.Printf("   - anti-spam check      : %s\n", format.ColoredBool(account.EnableAntiSpamCheck))
		// fmt.Printf("   - SMS unlimited credit : %s\n", format.ColoredBool(account.EnableSmsUnlimitedCredit))
		// fmt.Printf("   - OTP SMS              : %s\n", format.ColoredBool(account.EnableOtpSms))
		// fmt.Printf("   - OTP email            : %s\n", format.ColoredBool(account.EnableOtpEmail))

		fmt.Printf("%s", format.ToYAML(gateway))

	}
	return nil
}
