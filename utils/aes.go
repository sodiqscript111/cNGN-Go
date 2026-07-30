package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/theobiabo/cNGN-Go/envelope"
	cngnErr "github.com/theobiabo/cNGN-Go/error"
)

const ivLength = 16

func makeKey(key string) [32]byte {
	return sha256.Sum256([]byte(key))
}

func AESEncrypt(data, key string) (*envelope.EncryptedPayload, *cngnErr.Error) {
	derivedKey := makeKey(key)
	iv := make([]byte, ivLength)
	if _, err := rand.Read(iv); err != nil {
		return nil, cngnErr.NewCryptoError(fmt.Sprintf("failed to generate IV: %v", err))
	}

	block, err := aes.NewCipher(derivedKey[:])
	if err != nil {
		return nil, cngnErr.NewCryptoError(fmt.Sprintf("failed to create cipher: %v", err))
	}

	padded := pkcs7Pad([]byte(data), aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return &envelope.EncryptedPayload{
		Content: base64.StdEncoding.EncodeToString(ciphertext),
		IV:      base64.StdEncoding.EncodeToString(iv),
	}, nil
}

func AESDecrypt(payload *envelope.EncryptedPayload, key string) (string, *cngnErr.Error) {
	iv, err := base64.StdEncoding.DecodeString(payload.IV)
	if err != nil {
		return "", cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("invalid IV: %v", err))
	}
	if len(iv) != ivLength {
		return "", cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("IV must be %d bytes", ivLength))
	}

	encryptedContent, err := base64.StdEncoding.DecodeString(payload.Content)
	if err != nil {
		return "", cngnErr.NewInvalidEncryptedPayloadError(fmt.Sprintf("invalid encrypted content: %v", err))
	}

	derivedKey := makeKey(key)
	block, err := aes.NewCipher(derivedKey[:])
	if err != nil {
		return "", cngnErr.NewCryptoError(fmt.Sprintf("failed to create cipher: %v", err))
	}

	if len(encryptedContent) < aes.BlockSize {
		return "", cngnErr.NewInvalidEncryptedPayloadError("encrypted content too short")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedContent))
	mode.CryptBlocks(decrypted, encryptedContent)

	var padErr string
	decrypted, padErr = pkcs7Unpad(decrypted, aes.BlockSize)
	if padErr != "" {
		return "", cngnErr.NewCryptoError(fmt.Sprintf("AES decryption failed: %s", padErr))
	}

	return string(decrypted), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, string) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, "invalid padding length"
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, "invalid padding value"
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, "invalid padding"
		}
	}
	return data[:len(data)-padding], ""
}
