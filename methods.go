package cngn

import (
	"fmt"

	"github.com/sodiqscript111/cNGN-Go/envelope"
	cngnErr "github.com/sodiqscript111/cNGN-Go/error"
)

func (c *Client) Balance() (*envelope.Response[[]envelope.Balance], *cngnErr.Error) {
	return SendRequest[[]envelope.Balance](c, "GET", "/balance", nil, nil)
}

func (c *Client) Transactions(page, limit uint32) (*envelope.Response[envelope.Paginated[envelope.Transaction]], *cngnErr.Error) {
	return SendRequest[envelope.Paginated[envelope.Transaction]](c, "GET", "/transactions", map[string]string{
		"page":  fmtUint32(page),
		"limit": fmtUint32(limit),
	}, nil)
}

func (c *Client) Networks() (*envelope.Response[[]envelope.Network], *cngnErr.Error) {
	return SendRequest[[]envelope.Network](c, "GET", "/networks", nil, nil)
}

func (c *Client) Whitelist(request *WhitelistRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/updateBusiness", nil, request)
}

func (c *Client) Whitelisted() (*envelope.Response[[]envelope.WhitelistedAddress], *cngnErr.Error) {
	return SendRequest[[]envelope.WhitelistedAddress](c, "GET", "/whitelisted", nil, nil)
}

func (c *Client) VirtualAccounts() (*envelope.Response[[]envelope.VirtualAccount], *cngnErr.Error) {
	return SendRequest[[]envelope.VirtualAccount](c, "GET", "/virtual-account", nil, nil)
}

func (c *Client) CreateTemporaryVirtualAccount(request *CreateVirtualAccountRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/virtual-account/temporary", nil, request)
}

func (c *Client) VerifyAccount(request *VerifyAccountRequest) (*envelope.Response[envelope.AccountVerification], *cngnErr.Error) {
	return SendRequest[envelope.AccountVerification](c, "POST", "/account/verify", nil, request)
}

func (c *Client) Banks() (*envelope.Response[[]envelope.Bank], *cngnErr.Error) {
	return SendRequest[[]envelope.Bank](c, "GET", "/banks", nil, nil)
}

func (c *Client) UpdateBankAccount(request *UpdateBankAccountRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "PUT", "/bank-account", nil, request)
}

func (c *Client) SubmitWithdraw(request *WithdrawRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/withdraw", nil, request)
}

func (c *Client) VerifyWithdrawal(reference string) (*envelope.Response[envelope.Transaction], *cngnErr.Error) {
	return SendRequest[envelope.Transaction](c, "GET", fmt.Sprintf("/withdraw/verify/%s", reference), nil, nil)
}

func (c *Client) BridgeQuote(request *BridgeQuoteRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/swap-quote", nil, request)
}

func (c *Client) Bridge(request *BridgeRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/swap", nil, request)
}

func (c *Client) RedeemAsset(request *RedeemAssetRequest) (envelope.JsonResponse, *cngnErr.Error) {
	return sendJSON[envelope.JsonResponse](c, "POST", "/redeemAsset", nil, request)
}
