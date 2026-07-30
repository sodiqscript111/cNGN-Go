package utils

import (
	"testing"

	"github.com/sodiqscript111/cNGN-Go/envelope"
	cngnErr "github.com/sodiqscript111/cNGN-Go/error"
)

func TestAESEncryptDecryptRoundTrip(t *testing.T) {
	key := "test-key-12345678"
	plaintext := `{"amount":"100","currency":"NGN"}`

	payload, err := AESEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("AESEncrypt failed: %v", err)
	}

	if payload.Content == "" {
		t.Error("expected non-empty Content")
	}
	if payload.IV == "" {
		t.Error("expected non-empty IV")
	}

	decrypted, err := AESDecrypt(payload, key)
	if err != nil {
		t.Fatalf("AESDecrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("round-trip mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestAESEncryptProducesDifferentCiphertexts(t *testing.T) {
	key := "fixed-key"
	plaintext := "hello world"

	p1, _ := AESEncrypt(plaintext, key)
	p2, _ := AESEncrypt(plaintext, key)

	if p1.Content == p2.Content {
		t.Error("expected different ciphertexts due to random IV")
	}
	if p1.IV == p2.IV {
		t.Error("expected different IVs")
	}
}

func TestAESDecryptWrongKey(t *testing.T) {
	payload, _ := AESEncrypt("secret data", "correct-key")
	_, err := AESDecrypt(payload, "wrong-key")
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
	if err.Kind != 	cngnErr.TypeCrypto {
		t.Errorf("expected crypto error kind, got %v", err.Kind)
	}
}

func TestAESDecryptInvalidBase64(t *testing.T) {
	payload := &envelope.EncryptedPayload{
		Content: "!!!not-base64!!!",
		IV:      "dGVzdC1pdi0tMTYtYnl0ZXM=",
	}
	_, err := AESDecrypt(payload, "key")
	if err == nil {
		t.Fatal("expected error for invalid base64 content")
	}
}

func TestAESDecryptInvalidIVLength(t *testing.T) {
	payload := &envelope.EncryptedPayload{
		Content: base64Encode([]byte("some-ciphertext")),
		IV:      base64Encode([]byte("short")),
	}
	_, err := AESDecrypt(payload, "key")
	if err == nil {
		t.Fatal("expected error for short IV")
	}
}

func TestAESDecryptTooShortContent(t *testing.T) {
	payload := &envelope.EncryptedPayload{
		Content: base64Encode([]byte("x")),
		IV:      base64Encode([]byte("1234567890123456")),
	}
	_, err := AESDecrypt(payload, "key")
	if err == nil {
		t.Fatal("expected error for content shorter than block size")
	}
}

func TestAESDecryptInvalidIVEncoding(t *testing.T) {
	payload := &envelope.EncryptedPayload{
		Content: "dGVzdA==",
		IV:      "!!!not-base64!!!",
	}
	_, err := AESDecrypt(payload, "key")
	if err == nil {
		t.Fatal("expected error for invalid base64 IV")
	}
}

func TestPKCS7Pad(t *testing.T) {
	tests := []struct {
		input     []byte
		blockSize int
		wantLen   int
	}{
		{[]byte(""), 16, 16},
		{[]byte("hello"), 16, 16},
		{[]byte("1234567890123456"), 16, 32},
		{[]byte("x"), 8, 8},
	}

	for _, tt := range tests {
		padded := pkcs7Pad(tt.input, tt.blockSize)
		if len(padded) != tt.wantLen {
			t.Errorf("pkcs7Pad(%q, %d) len = %d, want %d", tt.input, tt.blockSize, len(padded), tt.wantLen)
		}
		if len(padded)%tt.blockSize != 0 {
			t.Errorf("padded length %d not multiple of block size %d", len(padded), tt.blockSize)
		}
		padByte := padded[len(padded)-1]
		for i := len(padded) - int(padByte); i < len(padded); i++ {
			if padded[i] != padByte {
				t.Errorf("padding byte mismatch at position %d", i)
			}
		}
	}
}

func TestPKCS7Unpad(t *testing.T) {
	tests := []struct {
		input     []byte
		blockSize int
		want      []byte
		wantErr   bool
	}{
		{append([]byte("hello"), bytesRepeat(11, 11)...), 16, []byte("hello"), false},
		{bytesRepeat(16, 16), 16, []byte{}, false},
		{nil, 16, nil, true},
		{[]byte{}, 16, nil, true},
		{[]byte("hello"), 16, nil, true},
		{append([]byte("data"), 0), 16, nil, true},
		{append([]byte("data"), 17), 16, nil, true},
	}

	for _, tt := range tests {
		result, err := pkcs7Unpad(tt.input, tt.blockSize)
		if tt.wantErr {
			if err == nil {
				t.Errorf("pkcs7Unpad(%v, %d) expected error", tt.input, tt.blockSize)
			}
			continue
		}
		if err != nil {
			t.Errorf("pkcs7Unpad(%v, %d) unexpected error: %v", tt.input, tt.blockSize, err)
			continue
		}
		if string(result) != string(tt.want) {
			t.Errorf("pkcs7Unpad = %q, want %q", string(result), string(tt.want))
		}
	}
}

func TestPKCS7RoundTrip(t *testing.T) {
	inputs := []string{"", "a", "hello", "1234567890123456", "12345678901234567890"}
	for _, in := range inputs {
		padded := pkcs7Pad([]byte(in), 16)
		unpadded, err := pkcs7Unpad(padded, 16)
		if err != nil {
			t.Errorf("unpad error for input %q: %v", in, err)
			continue
		}
		if string(unpadded) != in {
			t.Errorf("round-trip mismatch: %q vs %q", string(unpadded), in)
		}
	}
}

func bytesRepeat(b byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = b
	}
	return out
}

func base64Encode(src []byte) string {
	dst := make([]byte, base64EncodedLen(len(src)))
	base64EncodeBytes(dst, src)
	return string(dst)
}

func base64EncodedLen(n int) int {
	return (n + 2) / 3 * 4
}

func base64EncodeBytes(dst, src []byte) {
	const enc = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	for len(src) > 0 {
		var val int
		switch len(src) {
		default:
			val = int(src[0])<<16 | int(src[1])<<8 | int(src[2])
			dst[0] = enc[(val>>18)&0x3F]
			dst[1] = enc[(val>>12)&0x3F]
			dst[2] = enc[(val>>6)&0x3F]
			dst[3] = enc[val&0x3F]
			src = src[3:]
			dst = dst[4:]
		case 2:
			val = int(src[0])<<16 | int(src[1])<<8
			dst[0] = enc[(val>>18)&0x3F]
			dst[1] = enc[(val>>12)&0x3F]
			dst[2] = enc[(val>>6)&0x3F]
			dst[3] = '='
			src = src[2:]
			dst = dst[4:]
		case 1:
			val = int(src[0]) << 16
			dst[0] = enc[(val>>18)&0x3F]
			dst[1] = enc[(val>>12)&0x3F]
			dst[2] = '='
			dst[3] = '='
			src = src[1:]
			dst = dst[4:]
		}
	}
}
