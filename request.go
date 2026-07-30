package cngn

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

type RedeemAssetRequest struct {
	Amount        string `json:"amount"`
	Bank          string `json:"bank"`
	AccountNumber string `json:"accountNumber"`
	SaveDetails   bool   `json:"saveDetails"`
}
