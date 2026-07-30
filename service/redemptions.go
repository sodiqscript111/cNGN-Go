package service

import (
	"github.com/theobiabo/cNGN-Go/envelope"
	"github.com/theobiabo/cNGN-Go/error"
)

type Redemptions struct {
	send sendFunc
}

func NewRedemptions(send sendFunc) *Redemptions {
	return &Redemptions{send: send}
}

func (r *Redemptions) RedeemAsset(request *RedeemAssetRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](r.send, "POST", "/redeemAsset", nil, request)
}

type RedeemAssetRequest struct {
	Amount        string `json:"amount"`
	Bank          string `json:"bank"`
	AccountNumber string `json:"accountNumber"`
	SaveDetails   bool   `json:"saveDetails"`
}
