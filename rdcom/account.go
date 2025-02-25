package rdcom

import (
	"errors"
	"log/slog"
	"time"

	"github.com/dihedron/sms/pointer"
)

type AccountService struct {
	Service
}

type Account struct {
	Name                     string    `json:"name"`
	Code                     string    `json:"code"`
	Enabled                  bool      `json:"enabled"`
	Created                  time.Time `json:"created"`
	Parent                   string    `json:"parent"`
	ExpirationDate           time.Time `json:"expiration_date"`
	EnableEmailPreview       bool      `json:"enable_email_preview"`
	EnablePdfAttachments     bool      `json:"enable_pdf_attachments"`
	EnableAntiSpamCheck      bool      `json:"enable_anti_spam_check"`
	EnableSmsUnlimitedCredit bool      `json:"enable_sms_unlimited_credit"`
	EnableOtpSms             bool      `json:"enable_otp_sms"`
	EnableOtpEmail           bool      `json:"enable_otp_email"`
	Domains                  []string  `json:"domains"`
	UserPermissions          []struct {
		Account                         string `json:"account"`
		User                            string `json:"user"`
		CanViewLists                    bool   `json:"can_view_lists"`
		CanEditLists                    bool   `json:"can_edit_lists"`
		CanExportLists                  bool   `json:"can_export_lists"`
		CanViewCampaigns                bool   `json:"can_view_campaigns"`
		CanEditCampaigns                bool   `json:"can_edit_campaigns"`
		CanSendCampaigns                bool   `json:"can_send_campaigns"`
		CanViewTemplates                bool   `json:"can_view_templates"`
		CanEditTemplates                bool   `json:"can_edit_templates"`
		CanViewSmsCampaigns             bool   `json:"can_view_sms_campaigns"`
		CanEditSmsCampaigns             bool   `json:"can_edit_sms_campaigns"`
		CanSendSmsCampaigns             bool   `json:"can_send_sms_campaigns"`
		CanUseAutomation                bool   `json:"can_use_automation"`
		CanRequestSmsSenderRegistration bool   `json:"can_request_sms_sender_registration"`
		CanViewSmsAnalytics             bool   `json:"can_view_sms_analytics"`
		CanManageSmsAPIConfigurations   bool   `json:"can_manage_sms_api_configurations"`
		CanManageSubAccounts            bool   `json:"can_manage_sub_accounts"`
		CanSendTransactionalSms         bool   `json:"can_send_transactional_sms"`
		CanViewLandingPages             bool   `json:"can_view_landing_pages"`
		CanEditLandingPages             bool   `json:"can_edit_landing_pages"`
		CanViewStatsReadersAndClickmap  bool   `json:"can_view_stats_readers_and_clickmap"`
		CanViewStatsGeolocation         bool   `json:"can_view_stats_geolocation"`
		CanViewStatsComparison          bool   `json:"can_view_stats_comparison"`
		CanViewStatsDevicesAndTrend     bool   `json:"can_view_stats_devices_and_trend"`
		CanSendOtpSms                   bool   `json:"can_send_otp_sms"`
		CanValidateOtpSms               bool   `json:"can_validate_otp_sms"`
		CanViewOtpSms                   bool   `json:"can_view_otp_sms"`
		CanRevokeOtpSms                 bool   `json:"can_revoke_otp_sms"`
		CanManageOtpSms                 bool   `json:"can_manage_otp_sms"`
		CanSendOtpEmail                 bool   `json:"can_send_otp_email"`
		CanValidateOtpEmail             bool   `json:"can_validate_otp_email"`
		CanViewOtpEmail                 bool   `json:"can_view_otp_email"`
		CanRevokeOtpEmail               bool   `json:"can_revoke_otp_email"`
		CanManageOtpEmail               bool   `json:"can_manage_otp_email"`
	} `json:"user_permissions"`
	Infos struct {
		Name               string `json:"name"`
		Email              string `json:"email"`
		MainContactName    string `json:"main_contact_name"`
		MainContactSurname string `json:"main_contact_surname"`
		MainContactEmail   string `json:"main_contact_email"`
		MainContactCell    string `json:"main_contact_cell"`
		ReprName           string `json:"repr_name"`
		ReprSurname        string `json:"repr_surname"`
		ReprCallme         string `json:"repr_callme"`
		ReprEmail          string `json:"repr_email"`
		ReprFiscalCode     string `json:"repr_fiscal_code"`
		Company            string `json:"company"`
		CompanyType        string `json:"company_type"`
		Vat                string `json:"vat"`
		City               string `json:"city"`
		Address            string `json:"address"`
		ZipCode            string `json:"zip_code"`
		State              string `json:"state"`
		Country            string `json:"country"`
		Phone              string `json:"phone"`
		Website            string `json:"website"`
	} `json:"infos"`
	SuspensionState int     `json:"suspension_state"`
	SmsCredists     float64 `json:"sms_credists"`
	SenderAddress   string  `json:"sender_address"`
	Limits          struct {
		MaxRecipientsPerDay   int `json:"max_recipients_per_day"`
		MaxRecipientsPerMonth int `json:"max_recipients_per_month"`
		MaxRecipientsPerYear  int `json:"max_recipients_per_year"`
		MaxLists              int `json:"max_lists"`
	} `json:"limits"`
}

// List returns the list of accounts.
func (a *AccountService) List() ([]Account, error) {
	if a.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	options := &PaginatedListOptions{
		Options: Options{
			EntityPath: "/api/v2/accounts",
		},
		PageSize: pointer.To(100),
	}

	result, err := PaginatedList[Account](a.client, options)

	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return result, nil
}
