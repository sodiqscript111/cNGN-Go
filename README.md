# cNGN SDK

Go client for the cNGN API.

## Installation

```bash
go get github.com/sodiqscript111/cNGN-Go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/sodiqscript111/cNGN-Go"
)

func main() {
	client, err := cngn.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	balance, err := client.Balance()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", balance)
}
```

For local testing, the SDK can load `CNGN_KEY`, `CNGN_ENCRYPTION_KEY`, and
`CNGN_PRIVATE_KEY` from the process environment or a local `.env` file. The
private key should be the complete OpenSSH Ed25519 private key, including its
`BEGIN OPENSSH PRIVATE KEY` and `END OPENSSH PRIVATE KEY` lines.

You can also provide the values directly:

```go
client := cngn.New("API_KEY").WithSecurity("encryption_key", "openssh_private_key")
```

## Testing

```bash
go test ./...
```

All tests use a temporary local HTTP server and check every service method
without an API key, encryption key, quota, or IP whitelist.

## Services

### Wallet manager

Provides balances, transactions, supported networks, wallet whitelisting, and whitelisted addresses.

```go
balance, _ := client.Balance()
transactions, _ := client.Transactions(1, 20)
networks, _ := client.Networks()
```

### Deposit

Provides virtual accounts, temporary virtual accounts, bank lookup, account verification, and settlement bank updates.

```go
accounts, _ := client.VirtualAccounts()
banks, _ := client.Banks()
```

### Redemptions

Provides cNGN redemption to a bank account.

```go
resp, _ := client.RedeemAsset(&cngn.RedeemAssetRequest{
	Amount: "1000", Bank: "001", AccountNumber: "1234567890",
})
```

### Onchain transfer

Provides wallet withdrawals, withdrawal verification, bridge quotes, and bridging.

```go
resp, _ := client.SubmitWithdraw(&cngn.WithdrawRequest{
	Amount: "100", Address: "0x...", Network: "ETH",
})
status, _ := client.VerifyWithdrawal("ref-123")
```

## Authentication

The client sends the API key as a Bearer token. The API also requires a whitelisted source IP.

API failures are returned as `*error.Error`. Match on `Error.Kind` for
configuration, network, API, parse, and crypto errors. `Error.IsRetryable()`
is useful for network, `429`, service-availability, and `5xx` failures.

## Base URL

The default base URL is `https://api.cngn.co/v1/api`.

For testing or a custom deployment:

```go
client := cngn.WithBaseURL("API_KEY", "https://api.cngn.co/v1/api")
```

## Encryption

The SDK encrypts request bodies with AES-256-CBC. The key is SHA-256 hashed,
the IV is random for every request, and both encrypted values are base64
encoded in the `content` and `iv` fields expected by cNGN. Encrypted responses
are opened with the configured OpenSSH Ed25519 private key via NaCl box
(Curve25519 precomputation).
