package service

import (
	"github.com/theobiabo/cNGN-Go/envelope"
	"github.com/theobiabo/cNGN-Go/error"
)

type WalletManager struct {
	send sendFunc
}

func NewWalletManager(send sendFunc) *WalletManager {
	return &WalletManager{send: send}
}

func (w *WalletManager) Balance() (*envelope.Response[[]envelope.Balance], *error.Error) {
	return sendRequest[[]envelope.Balance](w.send, "GET", "/balance", nil, nil)
}

func (w *WalletManager) Transactions(page, limit uint32) (*envelope.Response[envelope.Paginated[envelope.Transaction]], *error.Error) {
	query := map[string]string{
		"page":  fmtUint32(page),
		"limit": fmtUint32(limit),
	}
	return sendRequest[envelope.Paginated[envelope.Transaction]](w.send, "GET", "/transactions", query, nil)
}

func (w *WalletManager) Networks() (*envelope.Response[[]envelope.Network], *error.Error) {
	return sendRequest[[]envelope.Network](w.send, "GET", "/networks", nil, nil)
}

func (w *WalletManager) Whitelist(request *WhitelistRequest) (envelope.JsonResponse, *error.Error) {
	return sendJSON[envelope.JsonResponse](w.send, "POST", "/updateBusiness", nil, request)
}

func (w *WalletManager) Whitelisted() (*envelope.Response[[]envelope.WhitelistedAddress], *error.Error) {
	return sendRequest[[]envelope.WhitelistedAddress](w.send, "GET", "/whitelisted", nil, nil)
}

type WhitelistRequest struct {
	WalletAddress *WalletAddress `json:"walletAddress,omitempty"`
	BankDetails   *BankDetails   `json:"bankDetails,omitempty"`
}

type WalletAddress struct {
	XbnAddress     *string `json:"xbnAddress,omitempty"`
	AtcAddress     *string `json:"atcAddress,omitempty"`
	BscAddress     *string `json:"bscAddress,omitempty"`
	EthAddress     *string `json:"ethAddress,omitempty"`
	BaseAddress    *string `json:"baseAddress,omitempty"`
	PolygonAddress *string `json:"polygonAddress,omitempty"`
	TronAddress    *string `json:"tronAddress,omitempty"`
	BantuUserID    *string `json:"bantuUserId,omitempty"`
}

type BankDetails struct {
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
