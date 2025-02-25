package rdcom

import (
	"errors"
	"log/slog"
)

type SMSGatewayService struct {
	Service
}

type SMSGateway struct {
	ID             int                `json:"id"`
	Label          map[string]string  `json:"label"`
	IsDefault      bool               `json:"is_default"`
	SenderReady    bool               `json:"sender_ready"`
	TwowayReady    bool               `json:"twoway_ready"`
	MoReady        bool               `json:"mo_ready"`
	RcsReady       bool               `json:"rcs_ready"`
	EnableSmsDlr   bool               `json:"enable_sms_dlr"`
	GatewayType    string             `json:"gateway_type"`
	GatewayTypeRaw int                `json:"gateway_type_raw"`
	Prices         map[string]float64 `json:"prices"`
}

// List returns the list of SMS gateways.
func (a *SMSGatewayService) List(account string) ([]SMSGateway, error) {
	if a.client.token == "" {
		slog.Error("invalid token")
		return nil, errors.New("invalid token")
	}

	options := &ListOptions{
		EntityPath: "/api/v2/{account}/cds/sms/",
		PathParams: map[string]string{
			"account": account,
		},
	}

	result, err := List[SMSGateway](a.client, options)

	if err != nil {
		slog.Error("error placing API call", "error", err)
		return nil, err
	}
	slog.Debug("API call success")
	return result, nil
}
