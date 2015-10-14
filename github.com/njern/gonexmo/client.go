package nexmo

import (
	"errors"
)

// Client encapsulates the Nexmo functions - must be created with
// NewClientFromAPI()
type Client struct {
	Account   *Account
	SMS       *SMS
	USSD      *USSD
	apiKey    string
	apiSecret string
	useOauth  bool
}

// NewClientFromAPI creates a new Client type with the
// provided API key / API secret.
func NewClientFromAPI(apiKey, apiSecret string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey can not be empty")
	} else if apiSecret == "" {
		return nil, errors.New("apiSecret can not be empty")
	}

	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		useOauth:  false,
	}

	c.Account = &Account{c}
	c.SMS = &SMS{c}
	c.USSD = &USSD{c}
	return c, nil
}
