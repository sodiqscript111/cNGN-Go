package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/nacl/box"

	cngnErr "github.com/theobiabo/cNGN-Go/error"
)

const (
	nonceLength     = 24
	publicKeyLength = 32
)

func Ed25519DecryptWithPrivateKey(privateKey, encryptedData string) (string, *cngnErr.Error) {
	ed25519Key, err := parseOpenSSHPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	curve25519Key := convertPrivateKey(ed25519Key)

	encryptedBytes, decErr := base64.StdEncoding.DecodeString(encryptedData)
	if decErr != nil {
		return "", cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("invalid response encoding: %v", decErr))
	}

	minimumLength := nonceLength + publicKeyLength
	if len(encryptedBytes) < minimumLength {
		return "", cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("encrypted response is too short; expected at least %d bytes", minimumLength))
	}

	nonceBytes := encryptedBytes[:nonceLength]
	ephemeralPubBytes := encryptedBytes[len(encryptedBytes)-publicKeyLength:]
	ciphertext := encryptedBytes[nonceLength : len(encryptedBytes)-publicKeyLength]

	var nonce [24]byte
	copy(nonce[:], nonceBytes)

	var ephemeralPub [32]byte
	copy(ephemeralPub[:], ephemeralPubBytes)

	var secretKey [32]byte
	copy(secretKey[:], curve25519Key[:])

	sharedKey := new([32]byte)
	box.Precompute(sharedKey, &ephemeralPub, &secretKey)
	decrypted, ok := box.OpenAfterPrecomputation(nil, ciphertext, &nonce, sharedKey)
	if !ok {
		return "", cngnErr.NewCryptoError("Ed25519 response decryption failed")
	}

	return string(decrypted), nil
}

func parseOpenSSHPrivateKey(privateKey string) ([64]byte, *cngnErr.Error) {
	header, rest, ok := strings.Cut(privateKey, "\n")
	if !ok {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError("private key must be an OpenSSH Ed25519 key")
	}

	header = strings.TrimSpace(header)
	if header != "-----BEGIN OPENSSH PRIVATE KEY-----" {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError("private key must be an OpenSSH Ed25519 key")
	}

	bodyBuilder := strings.Builder{}
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "-----END OPENSSH PRIVATE KEY-----" {
			break
		}
		bodyBuilder.WriteString(trimmed)
	}

	decoded, err := base64.StdEncoding.DecodeString(bodyBuilder.String())
	if err != nil {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("invalid OpenSSH private key: %v", err))
	}

	marker := []byte{0, 0, 0, 0x40}
	markerStart := -1
	for i := 0; i <= len(decoded)-len(marker); i++ {
		match := true
		for j := 0; j < len(marker); j++ {
			if decoded[i+j] != marker[j] {
				match = false
				break
			}
		}
		if match {
			markerStart = i
			break
		}
	}

	if markerStart < 0 {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError("Ed25519 key data was not found")
	}

	keyStart := markerStart + len(marker)
	keyEnd := keyStart + 64
	if keyEnd > len(decoded) {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError("OpenSSH private key is incomplete")
	}

	return toArray64(decoded[keyStart:keyEnd])
}

func convertPrivateKey(ed25519Key [64]byte) [32]byte {
	seed := ed25519Key[:32]
	hash := sha512.Sum512(seed)
	hash[0] &= 248
	hash[31] &= 127
	hash[31] |= 64
	return toArray32(hash[:32])
}

func toArray32(bytes []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], bytes)
	return arr
}

func toArray64(bytes []byte) ([64]byte, *cngnErr.Error) {
	if len(bytes) != 64 {
		return [64]byte{}, cngnErr.NewInvalidEncryptedPayloadError("Ed25519 private key must contain 64 bytes")
	}
	var arr [64]byte
	copy(arr[:], bytes)
	return arr, nil
}
