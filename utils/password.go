package utils

import(
	"crypto/sha256"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)


func Hash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return string(sum[:])
}

func Encrypt(plaintext string, key []byte) (string, error) {
	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating cipher block: %v", err)
	}

	// Create initialization vector
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("error creating IV: %v", err)
	}

	// Pad the plaintext
	padded := pad([]byte(plaintext))

	// Create ciphertext slice
	ciphertext := make([]byte, len(iv)+len(padded))
	copy(ciphertext, iv)

	// Create CBC encrypter
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	// Encode to base64 and return
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt takes a base64-encoded ciphertext and key, returns the original plaintext
func Decrypt(encodedCiphertext string, key []byte) (string, error) {
	// Decode base64
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("error decoding base64: %v", err)
	}

	// Check for minimum length
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating cipher block: %v", err)
	}

	// Extract IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad the result
	unpadded, err := unpad(ciphertext)
	if err != nil {
		return "", fmt.Errorf("error unpadding: %v", err)
	}

	return string(unpadded), nil
}

// pad implements PKCS7 padding
func pad(data []byte) []byte {
	padLen := aes.BlockSize - (len(data) % aes.BlockSize)
	padding := make([]byte, padLen)
	for i := range padding {
		padding[i] = byte(padLen)
	}
	return append(data, padding...)
}

// unpad removes PKCS7 padding
func unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	padLen := int(data[len(data)-1])
	if padLen > len(data) || padLen > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}

	// Verify padding
	for i := len(data) - padLen; i < len(data); i++ {
		if data[i] != byte(padLen) {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-padLen], nil
}