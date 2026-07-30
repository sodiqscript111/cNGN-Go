package envelope

type Response[T any] struct {
	Status  uint16 `json:"status"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type ErrorEnvelope struct {
	Status  interface{} `json:"status"`
	Message string      `json:"message"`
}

type Pagination struct {
	Count        uint32  `json:"count"`
	Pages        uint32  `json:"pages"`
	IsLastPage   bool    `json:"is_last_page"`
	NextPage     *uint32 `json:"next_page,omitempty"`
	PreviousPage *uint32 `json:"previous_page,omitempty"`
}

type Paginated[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Balance struct {
	AssetType string `json:"asset_type"`
	AssetCode string `json:"asset_code"`
	Balance   string `json:"balance"`
}

type Bank struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Code         string `json:"code"`
	Country      string `json:"country"`
	NibssBankCode string `json:"nibss_bank_code"`
}

type Network struct {
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	ChainID *string `json:"chain_id,omitempty"`
}

type Receiver struct {
	Address       *string `json:"address,omitempty"`
	Bank          *Bank   `json:"bank,omitempty"`
	AccountNumber *string `json:"account_number,omitempty"`
}

type Transaction struct {
	ID            string    `json:"id"`
	From          string    `json:"from"`
	Receiver      Receiver  `json:"receiver"`
	Amount        string    `json:"amount"`
	Description   string    `json:"description"`
	CreatedAt     string    `json:"created_at"`
	TrxRef        string    `json:"trx_ref"`
	TrxType       string    `json:"trx_type"`
	Network       string    `json:"network"`
	AssetType     string    `json:"asset_type"`
	AssetSymbol   string    `json:"asset_symbol"`
	BaseTrxHash   *string   `json:"base_trx_hash,omitempty"`
	ExtlTrxHash   *string   `json:"extl_trx_hash,omitempty"`
	ExplorerLink  *string   `json:"explorer_link,omitempty"`
	Status        string    `json:"status"`
	UpdatedAt     *string   `json:"updated_at,omitempty"`
}

type VirtualAccount struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	Bank          *Bank  `json:"bank,omitempty"`
	Provider      *string `json:"provider,omitempty"`
}

type AccountVerification struct {
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
}

type Withdrawal struct {
	TrxRef string  `json:"trx_ref"`
	Status *string `json:"status,omitempty"`
}

type WhitelistedAddress struct {
	Address string  `json:"address"`
	Network string  `json:"network"`
	Label   *string `json:"label,omitempty"`
}

type EncryptedPayload struct {
	Content string `json:"content"`
	IV      string `json:"iv"`
}

type JsonResponse = Response[interface{}]

type BalanceResponse = Response[[]Balance]
type TransactionResponse = Response[Transaction]
type PaginatedTransactionResponse = Response[Paginated[Transaction]]
type NetworkResponse = Response[[]Network]
type WhitelistedAddressResponse = Response[[]WhitelistedAddress]
type VirtualAccountResponse = Response[[]VirtualAccount]
type AccountVerificationResponse = Response[AccountVerification]
type BankResponse = Response[[]Bank]
