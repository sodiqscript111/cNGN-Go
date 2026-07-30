package utils

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

const testPrivateKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAAHOewbzQPJaAL74Phk6tysEIohzeuQGB/Y6WBBnpXMgAAAJBP07vTT9O7
0wAAAAtzc2gtZWQyNTUxOQAAACAAHOewbzQPJaAL74Phk6tysEIohzeuQGB/Y6WBBnpXMg
AAAED+u9nZ2tiy3IBD10P0mknxE47io8aa9ks9BL4ny4jILAAc57BvNA8loAvvg+GTq3Kw
QiiHN65AYH9jpYEGelcyAAAADHRlc3RAY25nbi1nbwE=
-----END OPENSSH PRIVATE KEY-----`

func TestParseOpenSSHPrivateKeyValid(t *testing.T) {
	key, err := parseOpenSSHPrivateKey(testPrivateKey)
	if err != nil {
		t.Fatalf("parseOpenSSHPrivateKey failed: %v", err)
	}
	if len(key) != 64 {
		t.Errorf("expected 64 bytes, got %d", len(key))
	}
}

func TestParseOpenSSHPrivateKeyMissingHeader(t *testing.T) {
	_, err := parseOpenSSHPrivateKey("not a key")
	if err == nil {
		t.Fatal("expected error for missing header")
	}
}

func TestParseOpenSSHPrivateKeyWrongHeader(t *testing.T) {
	invalid := "-----BEGIN RSA PRIVATE KEY-----\ndata\n-----END RSA PRIVATE KEY-----"
	_, err := parseOpenSSHPrivateKey(invalid)
	if err == nil {
		t.Fatal("expected error for wrong header")
	}
}

func TestParseOpenSSHPrivateKeyInvalidBase64(t *testing.T) {
	invalid := "-----BEGIN OPENSSH PRIVATE KEY-----\n!!!not-base64!!!\n-----END OPENSSH PRIVATE KEY-----"
	_, err := parseOpenSSHPrivateKey(invalid)
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestParseOpenSSHPrivateKeyTooShort(t *testing.T) {
	shortBody := base64.StdEncoding.EncodeToString([]byte("short"))
	invalid := "-----BEGIN OPENSSH PRIVATE KEY-----\n" + shortBody + "\n-----END OPENSSH PRIVATE KEY-----"
	_, err := parseOpenSSHPrivateKey(invalid)
	if err == nil {
		t.Fatal("expected error for too-short key data")
	}
}
func TestEd25519DecryptWithPrivateKey(t *testing.T) {
	ed25519Key, err := parseOpenSSHPrivateKey(testPrivateKey)
	if err != nil {
		t.Fatalf("parseOpenSSHPrivateKey failed: %v", err)
	}

	curve25519Priv := convertPrivateKey(ed25519Key)

	// Derive the Curve25519 public key from the private key
	curve25519PubBytes, dhErr := curve25519.X25519(curve25519Priv[:], curve25519.Basepoint)
	if dhErr != nil {
		t.Fatal(dhErr)
	}
	var curve25519Pub [32]byte
	copy(curve25519Pub[:], curve25519PubBytes)

	plaintext := `{"data":"hello-world"}`
	var nonce [24]byte
	if _, randErr := rand.Read(nonce[:]); randErr != nil {
		t.Fatal(randErr)
	}

	ephemeralPub, ephemeralSecret, genErr := box.GenerateKey(rand.Reader)
	if genErr != nil {
		t.Fatal(genErr)
	}

	sharedKey := new([32]byte)
	box.Precompute(sharedKey, &curve25519Pub, ephemeralSecret)
	encrypted := box.SealAfterPrecomputation(nil, []byte(plaintext), &nonce, sharedKey)

	fullPayload := make([]byte, len(nonce)+len(encrypted)+len(ephemeralPub))
	copy(fullPayload[:], nonce[:])
	copy(fullPayload[len(nonce):], encrypted)
	copy(fullPayload[len(nonce)+len(encrypted):], ephemeralPub[:])

	encodedPayload := base64.StdEncoding.EncodeToString(fullPayload)

	result, decryptErr := Ed25519DecryptWithPrivateKey(testPrivateKey, encodedPayload)
	if decryptErr != nil {
		t.Fatalf("Ed25519DecryptWithPrivateKey failed: %v", decryptErr)
	}
	if result != plaintext {
		t.Errorf("decrypted text mismatch: got %q, want %q", result, plaintext)
	}
}

func TestEd25519DecryptInvalidKey(t *testing.T) {
	_, err := Ed25519DecryptWithPrivateKey("not-a-key", "dGVzdA==")
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
}

func TestEd25519DecryptInvalidBase64(t *testing.T) {
	_, err := Ed25519DecryptWithPrivateKey(testPrivateKey, "!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64 payload")
	}
}

func TestEd25519DecryptTooShort(t *testing.T) {
	shortPayload := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err := Ed25519DecryptWithPrivateKey(testPrivateKey, shortPayload)
	if err == nil {
		t.Fatal("expected error for too-short payload")
	}
}

func TestConvertPrivateKeyClamping(t *testing.T) {
	var seed [64]byte
	for i := 0; i < 64; i++ {
		seed[i] = byte(i)
	}

	result := convertPrivateKey(seed)

	if result[0]&248 != result[0] {
		t.Error("first byte should have lower 3 bits cleared")
	}
	if result[31]&127 != result[31] {
		t.Error("last byte should have top bit cleared")
	}
	if result[31]|64 != result[31] {
		t.Error("last byte should have bit 6 set")
	}
}

func TestToArray64Valid(t *testing.T) {
	input := make([]byte, 64)
	for i := range input {
		input[i] = byte(i)
	}

	result, err := toArray64(input)
	if err != nil {
		t.Fatalf("toArray64 failed: %v", err)
	}
	if result[0] != 0 || result[63] != 63 {
		t.Errorf("unexpected content: first=%d last=%d", result[0], result[63])
	}
}

func TestToArray64WrongLength(t *testing.T) {
	_, err := toArray64([]byte{1, 2, 3})
	if err == nil {
		t.Fatal("expected error for wrong length")
	}
}
