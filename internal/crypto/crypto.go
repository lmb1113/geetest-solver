package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
)

func AesEncrypt(text []byte, secKey []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(secKey)

	if err != nil {
		return nil, err
	}

	paddedText := pkcs7Padding(text, block.BlockSize())
	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	return ciphertext, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func RsaEncrypt(plainText string, publicKey *rsa.PublicKey) (string, error) {
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptedBytes), nil
}

func RandomHex(length int) (string, error) {
	if length%2 != 0 {
		return "", fmt.Errorf("length must be an even number")
	}

	rndBytes := make([]byte, length/2)

	_, err := rand.Read(rndBytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rndBytes), nil
}
