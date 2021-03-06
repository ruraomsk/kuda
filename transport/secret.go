package transport

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func decode(message string) ([]byte, error) {
	inbuffer, _ := base64.StdEncoding.DecodeString(message)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(inbuffer) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := inbuffer[:aes.BlockSize]
	cbc := cipher.NewCBCDecrypter(block, iv)
	inbuffer = inbuffer[aes.BlockSize:]
	if len(inbuffer)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("длина буфера %d не кратна blocksize %d", len(inbuffer), aes.BlockSize)
	}
	cbc.CryptBlocks(inbuffer, inbuffer)
	return []byte(inbuffer), nil
}
func code(message []byte) (string, error) {
	for len(message)%aes.BlockSize != 0 {
		message = append(message, ' ')
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	chipText := make([]byte, aes.BlockSize+len(message))
	iv := chipText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(chipText[aes.BlockSize:], message)

	return base64.StdEncoding.EncodeToString(chipText), nil
}
