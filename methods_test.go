package cngn

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() (*httptest.Server, *Client) {
	mux := http.NewServeMux()

	mux.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": []map[string]string{
				{"asset_type": "NGN", "asset_code": "NGN", "balance": "1000"},
			},
		})
	})

	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]any{
				"data": []map[string]string{},
				"pagination": map[string]any{
					"count": 0, "pages": 1, "is_last_page": true,
				},
			},
		})
	})

	mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": []map[string]string{
				{"name": "Ethereum", "code": "ETH"},
			},
		})
	})

	mux.HandleFunc("/updateBusiness", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "success"},
		})
	})

	mux.HandleFunc("/whitelisted", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": []map[string]string{},
		})
	})

	mux.HandleFunc("/virtual-account/temporary", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "created"},
		})
	})

	mux.HandleFunc("/virtual-account", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": []map[string]string{
				{"account_number": "1234567890", "account_name": "Test"},
			},
		})
	})

	mux.HandleFunc("/account/verify", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{
				"account_name": "John", "account_number": "12345", "bank_code": "001",
			},
		})
	})

	mux.HandleFunc("/banks", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": []map[string]string{
				{"name": "Test Bank", "code": "001", "slug": "test-bank", "country": "NG", "nibss_bank_code": "111"},
			},
		})
	})

	mux.HandleFunc("/bank-account", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "updated"},
		})
	})

	mux.HandleFunc("/withdraw", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"trx_ref": "ref-123"},
		})
	})

	mux.HandleFunc("/withdraw/verify/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{
				"id": "tx-1", "status": "completed", "trx_ref": "ref-123",
			},
		})
	})

	mux.HandleFunc("/swap-quote", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "quoted"},
		})
	})

	mux.HandleFunc("/swap", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "bridged"},
		})
	})

	mux.HandleFunc("/redeemAsset", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"status": 200, "message": "OK",
			"data": map[string]string{"status": "redeemed"},
		})
	})

	srv := httptest.NewServer(mux)
	return srv, WithBaseURL("test-token", srv.URL)
}

func TestBalance(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Balance()
	if err != nil {
		t.Fatalf("Balance failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestTransactions(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Transactions(1, 20)
	if err != nil {
		t.Fatalf("Transactions failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestNetworks(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Networks()
	if err != nil {
		t.Fatalf("Networks failed: %v", err)
	}
	if len(*resp.Data) == 0 {
		t.Error("expected at least one network")
	}
}

func TestWhitelist(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Whitelist(&WhitelistRequest{
		BankDetails: &BankDetails{
			BankName: "Test", BankAccountName: "Owner", BankAccountNumber: "12345",
		},
	})
	if err != nil {
		t.Fatalf("Whitelist failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestWhitelisted(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Whitelisted()
	if err != nil {
		t.Fatalf("Whitelisted failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestVirtualAccounts(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.VirtualAccounts()
	if err != nil {
		t.Fatalf("VirtualAccounts failed: %v", err)
	}
	if len(*resp.Data) == 0 {
		t.Error("expected at least one virtual account")
	}
}

func TestCreateTemporaryVirtualAccount(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.CreateTemporaryVirtualAccount(&CreateVirtualAccountRequest{Provider: "test"})
	if err != nil {
		t.Fatalf("CreateTemporaryVirtualAccount failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestVerifyAccount(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.VerifyAccount(&VerifyAccountRequest{BankCode: "001", AccountNumber: "12345"})
	if err != nil {
		t.Fatalf("VerifyAccount failed: %v", err)
	}
	if resp.Data.AccountName != "John" {
		t.Errorf("expected 'John', got %q", resp.Data.AccountName)
	}
}

func TestBanks(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Banks()
	if err != nil {
		t.Fatalf("Banks failed: %v", err)
	}
	if len(*resp.Data) == 0 {
		t.Error("expected at least one bank")
	}
}

func TestUpdateBankAccount(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.UpdateBankAccount(&UpdateBankAccountRequest{
		BankCode: "001", AccountNumber: "12345", AccountName: "John",
	})
	if err != nil {
		t.Fatalf("UpdateBankAccount failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestSubmitWithdraw(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.SubmitWithdraw(&WithdrawRequest{
		Amount: "100", Address: "0xabc", Network: "ETH",
	})
	if err != nil {
		t.Fatalf("SubmitWithdraw failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestVerifyWithdrawal(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.VerifyWithdrawal("ref-123")
	if err != nil {
		t.Fatalf("VerifyWithdrawal failed: %v", err)
	}
	if resp.Data.Status != "completed" {
		t.Errorf("expected 'completed', got %q", resp.Data.Status)
	}
}

func TestBridgeQuote(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.BridgeQuote(&BridgeQuoteRequest{
		DestinationNetwork: "polygon",
		DestinationAddress: "0xabc",
		OriginNetwork:      "ethereum",
		Amount:             "100",
	})
	if err != nil {
		t.Fatalf("BridgeQuote failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestBridge(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	url := "https://callback.example.com"
	resp, err := client.Bridge(&BridgeRequest{
		DestinationNetwork: "polygon",
		DestinationAddress: "0xabc",
		OriginNetwork:      "ethereum",
		Amount:             "100",
		CallbackURL:        &url,
	})
	if err != nil {
		t.Fatalf("Bridge failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestRedeemAsset(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.RedeemAsset(&RedeemAssetRequest{
		Amount: "500", Bank: "001", AccountNumber: "12345", SaveDetails: true,
	})
	if err != nil {
		t.Fatalf("RedeemAsset failed: %v", err)
	}
	if resp.Status != 200 {
		t.Errorf("expected status 200, got %d", resp.Status)
	}
}

func TestVerifyAccountWithData(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.VerifyAccount(&VerifyAccountRequest{BankCode: "001", AccountNumber: "12345"})
	if err != nil {
		t.Fatalf("VerifyAccount failed: %v", err)
	}
	if resp.Data.AccountName != "John" {
		t.Errorf("expected account name 'John', got %q", resp.Data.AccountName)
	}
	if resp.Data.AccountNumber != "12345" {
		t.Errorf("expected account number '12345', got %q", resp.Data.AccountNumber)
	}
	if resp.Data.BankCode != "001" {
		t.Errorf("expected bank code '001', got %q", resp.Data.BankCode)
	}
}

func TestTransactionsWithPagination(t *testing.T) {
	srv, client := setupTestServer()
	defer srv.Close()

	resp, err := client.Transactions(1, 20)
	if err != nil {
		t.Fatalf("Transactions failed: %v", err)
	}
	if !resp.Data.Pagination.IsLastPage {
		t.Error("expected isLastPage to be true")
	}
}


