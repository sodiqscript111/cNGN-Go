package service

import (
	"github.com/theobiabo/cNGN-Go/envelope"
	"github.com/theobiabo/cNGN-Go/error"
)

type Deposit struct {
	send sendFunc
}

func NewDeposit(send sendFunc) *Deposit {
	return &Deposit{send: send}
}

func (d *Deposit) VirtualAccounts() (*envelope.Response[[]envelope.VirtualAccount], *error.Error) {
	return sendRequest[[]envelope.VirtualAccount](d.send, "GET", "/virtual-account", nil, nil)
}

func (d *Deposit) CreateTemporaryVirtualAccount(request *CreateVirtualAccountRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](d.send, "POST", "/virtual-account/temporary", nil, request)
}

func (d *Deposit) VerifyAccount(request *VerifyAccountRequest) (*envelope.Response[envelope.AccountVerification], *error.Error) {
	return sendRequest[envelope.AccountVerification](d.send, "POST", "/account/verify", nil, request)
}

func (d *Deposit) Banks() (*envelope.Response[[]envelope.Bank], *error.Error) {
	return sendRequest[[]envelope.Bank](d.send, "GET", "/banks", nil, nil)
}

func (d *Deposit) UpdateBankAccount(request *UpdateBankAccountRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](d.send, "PUT", "/bank-account", nil, request)
}

type CreateVirtualAccountRequest struct {
	Provider string `json:"provider,omitempty"`
}

type VerifyAccountRequest struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
}

type UpdateBankAccountRequest struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}
