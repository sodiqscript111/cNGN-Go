package cngn

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sodiqscript111/cNGN-Go/envelope"
	cngnErr "github.com/sodiqscript111/cNGN-Go/error"
)

func TestNew(t *testing.T) {
	c := New("token-123")
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if c.authToken != "token-123" {
		t.Errorf("expected token 'token-123', got %q", c.authToken)
	}
	if c.baseURL != BaseURL {
		t.Errorf("expected baseURL %q, got %q", BaseURL, c.baseURL)
	}
	if c.http == nil {
		t.Error("expected non-nil http client")
	}
}

func TestWithBaseURL(t *testing.T) {
	c := WithBaseURL("tok", "https://example.com/api/")
	if c.baseURL != "https://example.com/api" {
		t.Errorf("expected trimmed baseURL, got %q", c.baseURL)
	}
}

func TestWithSecurity(t *testing.T) {
	c := New("tok").WithSecurity("enc-key", "priv-key")
	if !c.hasSecurity {
		t.Error("expected hasSecurity to be true")
	}
	if c.encryptionKey != "enc-key" {
		t.Errorf("expected encryptionKey 'enc-key', got %q", c.encryptionKey)
	}
	if c.privateKey != "priv-key" {
		t.Errorf("expected privateKey 'priv-key', got %q", c.privateKey)
	}
}

func clearEnvForTest() {
	os.Unsetenv("CNGN_KEY")
	os.Unsetenv("CNGN_ENCRYPTION_KEY")
	os.Unsetenv("CNGN_PRIVATE_KEY")
	os.Unsetenv("CNGN_SSH_PRIVATE_KEY")
}

func TestFromEnvSuccess(t *testing.T) {
	clearEnvForTest()
	dir := t.TempDir()
	envFile := dir + "\\.env"
	if err := os.WriteFile(envFile, []byte(
		"CNGN_KEY=test-key\nCNGN_ENCRYPTION_KEY=enc-key\nCNGN_PRIVATE_KEY=priv-key\n",
	), 0644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	client, err := FromEnv()
	if err != nil {
		t.Fatalf("FromEnv failed: %v", err)
	}
	if client.authToken != "test-key" {
		t.Errorf("expected authToken 'test-key', got %q", client.authToken)
	}
	if !client.hasSecurity {
		t.Error("expected hasSecurity to be true")
	}
}

func TestFromEnvMissingKey(t *testing.T) {
	clearEnvForTest()
	dir := t.TempDir()
	envFile := dir + "\\.env"
	if err := os.WriteFile(envFile, []byte("CNGN_ENCRYPTION_KEY=ek\nCNGN_PRIVATE_KEY=pk\n"), 0644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for missing CNGN_KEY")
	}
}

func TestFromEnvMissingEncryptionKey(t *testing.T) {
	clearEnvForTest()
	dir := t.TempDir()
	envFile := dir + "\\.env"
	if err := os.WriteFile(envFile, []byte("CNGN_KEY=tk\nCNGN_PRIVATE_KEY=pk\n"), 0644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for missing encryption key")
	}
}

func TestFromEnvMissingPrivateKey(t *testing.T) {
	clearEnvForTest()
	dir := t.TempDir()
	envFile := dir + "\\.env"
	if err := os.WriteFile(envFile, []byte("CNGN_KEY=tk\nCNGN_ENCRYPTION_KEY=ek\n"), 0644); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	_, err := FromEnv()
	if err == nil {
		t.Fatal("expected error for missing private key")
	}
}

func TestSendSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" {
			t.Errorf("bad auth: %q", r.Header.Get("Authorization"))
		}
		w.Write([]byte(`{"status":200,"message":"OK","data":{"balance":"100"}}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	result := &envelope.Response[map[string]string]{}
	err := c.Send("GET", "/balance", nil, nil, result)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if result.Status != 200 {
		t.Errorf("expected status 200, got %d", result.Status)
	}
	if (*result.Data)["balance"] != "100" {
		t.Errorf("expected balance '100', got %q", (*result.Data)["balance"])
	}
}

func TestSendErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"status":401,"message":"Invalid api key"}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	err := c.Send("GET", "/test", nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	if err.Kind != cngnErr.TypeApi {
		t.Errorf("expected TypeApi, got %v", err.Kind)
	}
	if err.ApiKind != cngnErr.InvalidApiKey {
		t.Errorf("expected InvalidApiKey, got %v", err.ApiKind)
	}
}

func TestSendQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1, got %q", r.URL.Query().Get("page"))
		}
		w.Write([]byte(`{"status":200,"message":"OK","data":[]}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	err := c.Send("GET", "/items", map[string]string{"page": "1"}, nil, nil)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestSendRequestBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json, got %q", r.Header.Get("Content-Type"))
		}
		w.Write([]byte(`{"status":200,"message":"OK","data":{}}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	err := c.Send("POST", "/submit", nil, map[string]string{"key": "val"}, nil)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
}

func TestSendNetworkError(t *testing.T) {
	c := WithBaseURL("tok", "http://127.0.0.1:1")
	err := c.Send("GET", "/test", nil, nil, nil)
	if err == nil {
		t.Fatal("expected network error")
	}
	if err.Kind != cngnErr.TypeNetwork {
		t.Errorf("expected TypeNetwork, got %v", err.Kind)
	}
}

func TestSendRequestGeneric(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":200,"message":"OK","data":{"balance":"50"}}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	resp, err := SendRequest[map[string]string](c, "GET", "/balance", nil, nil)
	if err != nil {
		t.Fatalf("SendRequest failed: %v", err)
	}
	if (*resp.Data)["balance"] != "50" {
		t.Errorf("expected balance '50', got %q", (*resp.Data)["balance"])
	}
}

func TestSendRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"status":404,"message":"Not found"}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL)
	_, err := SendRequest[map[string]string](c, "GET", "/nonexistent", nil, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSendWithEncryption(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// With encryption: body should be {"content":"...","iv":"..."}
		w.Write([]byte(`{"status":200,"message":"OK","data":{}}`))
	}))
	defer srv.Close()

	c := WithBaseURL("tok", srv.URL).WithSecurity("test-encryption-key", "test-private-key")
	err := c.Send("POST", "/secure", nil, map[string]string{"amount": "100"}, nil)
	if err != nil {
		t.Fatalf("encrypted Send failed: %v", err)
	}
}

func TestPrepareRequestBodyEncryption(t *testing.T) {
	c := New("tok").WithSecurity("enc-key", "priv-key")

	body, err := c.prepareRequestBody(map[string]string{"msg": "hello"})
	if err != nil {
		t.Fatalf("prepareRequestBody failed: %v", err)
	}

	if !strings.Contains(string(body), `"content"`) {
		t.Error("expected encrypted body to contain content field")
	}
	if !strings.Contains(string(body), `"iv"`) {
		t.Error("expected encrypted body to contain iv field")
	}
}

func TestPrepareRequestBodyNil(t *testing.T) {
	c := New("tok")
	body, err := c.prepareRequestBody(nil)
	if err != nil {
		t.Fatalf("prepareRequestBody(nil) failed: %v", err)
	}
	if body != nil {
		t.Errorf("expected nil body, got %v", body)
	}
}
