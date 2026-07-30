package error_test

import (
	"errors"
	"testing"

	cngnErr "github.com/sodiqscript111/cNGN-Go/error"
)

func TestConfigurationError(t *testing.T) {
	err := cngnErr.NewConfigurationError("missing key")
	if err.Kind != cngnErr.TypeConfiguration {
		t.Errorf("expected TypeConfiguration, got %v", err.Kind)
	}
	if err.Message != "configuration error: missing key" {
		t.Errorf("unexpected message: %s", err.Message)
	}
	if err.Err != nil {
		t.Errorf("expected nil Err")
	}
	if err.IsRetryable() {
		t.Error("configuration errors should not be retryable")
	}
}

func TestNetworkError(t *testing.T) {
	cause := errors.New("connection refused")
	err := cngnErr.NewNetworkError(cause)
	if err.Kind != cngnErr.TypeNetwork {
		t.Errorf("expected TypeNetwork, got %v", err.Kind)
	}
	if !err.IsRetryable() {
		t.Error("network errors should be retryable")
	}
	if !errors.Is(err, cause) {
		t.Error("NetworkError should wrap the original error")
	}
}

func TestParseError(t *testing.T) {
	cause := errors.New("invalid character")
	err := cngnErr.NewParseError(cause)
	if err.Kind != cngnErr.TypeParse {
		t.Errorf("expected TypeParse, got %v", err.Kind)
	}
	if err.IsRetryable() {
		t.Error("parse errors should not be retryable")
	}
}

func TestCryptoError(t *testing.T) {
	err := cngnErr.NewCryptoError("decryption failed")
	if err.Kind != cngnErr.TypeCrypto {
		t.Errorf("expected TypeCrypto, got %v", err.Kind)
	}
	if err.IsRetryable() {
		t.Error("crypto errors should not be retryable")
	}
}

func TestInvalidEncryptedPayloadError(t *testing.T) {
	err := cngnErr.NewInvalidEncryptedPayloadError("too short")
	if err.Kind != cngnErr.TypeInvalidEncryptedPayload {
		t.Errorf("expected TypeInvalidEncryptedPayload, got %v", err.Kind)
	}
}

func TestApiError(t *testing.T) {
	err := cngnErr.NewApiError(401, cngnErr.InvalidApiKey, "", "bad key")
	if err.Kind != cngnErr.TypeApi {
		t.Errorf("expected TypeApi, got %v", err.Kind)
	}
	if err.Status != 401 {
		t.Errorf("expected status 401, got %d", err.Status)
	}
	if err.ApiKind != cngnErr.InvalidApiKey {
		t.Errorf("expected InvalidApiKey, got %v", err.ApiKind)
	}
}

func TestApiErrorIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		status   uint16
		kind     cngnErr.ApiErrorKind
		want     bool
	}{
		{"server error", 500, cngnErr.Unknown, true},
		{"too many requests", 429, cngnErr.TooManyRequests, true},
		{"service unavailable", 503, cngnErr.ServiceUnavailable, true},
		{"client error", 400, cngnErr.Validation, false},
		{"unauthorized", 401, cngnErr.InvalidApiKey, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cngnErr.NewApiError(tt.status, tt.kind, "", "msg")
			if err.IsRetryable() != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", err.IsRetryable(), tt.want)
			}
		})
	}
}

func TestClassifyApiError(t *testing.T) {
	tests := []struct {
		message string
		kind    cngnErr.ApiErrorKind
	}{
		{"No token provided", cngnErr.NoTokenProvided},
		{"Invalid token prefix", cngnErr.InvalidTokenPrefix},
		{"Invalid api key", cngnErr.InvalidApiKey},
		{"Merchant not found", cngnErr.MerchantNotFound},
		{"No Test SSH Key found", cngnErr.NoTestSshKeyFound},
		{"No Live SSH Key found", cngnErr.NoLiveSshKeyFound},
		{"IP address not whitelisted", cngnErr.IpAddressNotWhitelisted},
		{"Could not determine client IP address", cngnErr.CouldNotDetermineClientIpAddress},
		{"Permission denied", cngnErr.PermissionDenied},
		{"Missing encryption data, key, or IV", cngnErr.MissingEncryptionDataKeyOrIv},
		{"Decryption failed", cngnErr.DecryptionFailed},
		{"Too many requests. Please try again later.", cngnErr.TooManyRequests},
		{"Service is currently unavailable. Please try again later.", cngnErr.ServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			kind, _, _ := cngnErr.ClassifyApiError(400, tt.message)
			if kind != tt.kind {
				t.Errorf("ClassifyApiError(400, %q) kind = %v, want %v", tt.message, kind, tt.kind)
			}
		})
	}
}

func TestClassifyBadRequestValidation(t *testing.T) {
	kind, field, msg := cngnErr.ClassifyApiError(400, "email: invalid format")
	if kind != cngnErr.Validation {
		t.Errorf("expected Validation, got %v", kind)
	}
	if field != "email" {
		t.Errorf("expected field 'email', got %q", field)
	}
	if msg != " invalid format" {
		t.Errorf("expected message ' invalid format', got %q", msg)
	}
}

func TestClassifyBadRequestBusinessLogic(t *testing.T) {
	kind, _, msg := cngnErr.ClassifyApiError(400, "Insufficient funds")
	if kind != cngnErr.BusinessLogic {
		t.Errorf("expected BusinessLogic, got %v", kind)
	}
	if msg != "Insufficient funds" {
		t.Errorf("expected message 'Insufficient funds', got %q", msg)
	}
}

func TestClassifyApiErrorUnknown(t *testing.T) {
	kind, _, msg := cngnErr.ClassifyApiError(500, "internal error")
	if kind != cngnErr.Unknown {
		t.Errorf("expected Unknown, got %v", kind)
	}
	if msg != "internal error" {
		t.Errorf("expected message 'internal error', got %q", msg)
	}
}

func TestApiErrorKindString(t *testing.T) {
	if cngnErr.NoTokenProvided.String() != "no token provided" {
		t.Errorf("unexpected string: %s", cngnErr.NoTokenProvided.String())
	}
	if cngnErr.Unknown.String() != "unknown API error" {
		t.Errorf("unexpected string: %s", cngnErr.Unknown.String())
	}
}

func TestErrorImplementsError(t *testing.T) {
	var _ error = cngnErr.NewConfigurationError("test")
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := cngnErr.NewNetworkError(cause)
	if !errors.Is(err, cause) {
		t.Error("errors.Is should find the wrapped cause")
	}
}

func TestErrorString(t *testing.T) {
	err := cngnErr.NewConfigurationError("test")
	if err.Error() != "configuration error: test" {
		t.Errorf("unexpected Error(): %s", err.Error())
	}
}
