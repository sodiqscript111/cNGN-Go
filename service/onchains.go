package service

import (
	"fmt"

	"github.com/theobiabo/cNGN-Go/envelope"
	"github.com/theobiabo/cNGN-Go/error"
)

type OnchainTransfer struct {
	send sendFunc
}

func NewOnchainTransfer(send sendFunc) *OnchainTransfer {
	return &OnchainTransfer{send: send}
}

func (o *OnchainTransfer) Withdraw(request *WithdrawRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](o.send, "POST", "/withdraw", nil, request)
}

func (o *OnchainTransfer) VerifyWithdrawal(reference string) (*envelope.Response[envelope.Transaction], *error.Error) {
	return sendRequest[envelope.Transaction](o.send, "GET", fmt.Sprintf("/withdraw/verify/%s", reference), nil, nil)
}

func (o *OnchainTransfer) BridgeQuote(request *BridgeQuoteRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](o.send, "POST", "/swap-quote", nil, request)
}

func (o *OnchainTransfer) Bridge(request *BridgeRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](o.send, "POST", "/swap", nil, request)
}

type WithdrawRequest struct {
	Amount            string `json:"amount"`
	Address           string `json:"address"`
	Network           string `json:"network"`
	ShouldSaveAddress bool   `json:"shouldSaveAddress"`
}

type BridgeQuoteRequest struct {
	DestinationNetwork string `json:"destinationNetwork"`
	DestinationAddress string `json:"destinationAddress"`
	OriginNetwork      string `json:"originNetwork"`
	Amount             string `json:"amount"`
}

type BridgeRequest struct {
	DestinationNetwork string  `json:"destinationNetwork"`
	DestinationAddress string  `json:"destinationAddress"`
	OriginNetwork      string  `json:"originNetwork"`
	Amount             string  `json:"amount"`
	CallbackURL        *string `json:"callbackUrl,omitempty"`
	SenderAddress      *string `json:"senderAddress,omitempty"`
}
