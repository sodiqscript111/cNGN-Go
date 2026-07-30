package cngn

import (
	"github.com/theobiabo/cNGN-Go/error"
	"github.com/theobiabo/cNGN-Go/service"
)

func (c *Client) WalletManager() *service.WalletManager {
	return service.NewWalletManager(c.Send)
}

func (c *Client) Deposit() *service.Deposit {
	return service.NewDeposit(c.Send)
}

func (c *Client) OnchainTransfer() *service.OnchainTransfer {
	return service.NewOnchainTransfer(c.Send)
}

func (c *Client) Redemptions() *service.Redemptions {
	return service.NewRedemptions(c.Send)
}

var _ = (*Client).WalletManager
var _ = (*Client).Deposit
var _ = (*Client).OnchainTransfer
var _ = (*Client).Redemptions

func NewClient(authToken string) *Client {
	return New(authToken)
}

func ClientFromEnv() (*Client, *error.Error) {
	return FromEnv()
}
