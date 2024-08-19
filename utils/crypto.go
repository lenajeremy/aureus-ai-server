package utils

import (
	"bytes"
	"code-review/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

// Encrypt encrypts plainText using AES in CBC mode.
func Encrypt(plainText string) string {
	key := mdHash(config.GetEnv("GITHUB_SECRET"))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}

	// Generate an IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Panicln(err)
	}

	// Pad the plaintext to a multiple of the block size
	plainTextBytes := pad([]byte(plainText), aes.BlockSize)

	ciphertext := make([]byte, len(iv)+len(plainTextBytes))
	copy(ciphertext[:aes.BlockSize], iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainTextBytes)

	return hex.EncodeToString(ciphertext)
}

// Decrypt decrypts cipheredText using AES in CBC mode.
func Decrypt(cipheredText string) string {
	ciphertext, err := hex.DecodeString(cipheredText)
	if err != nil {
		log.Panicln(err)
	}

	key := mdHash(config.GetEnv("GITHUB_SECRET"))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(err)
	}

	if len(ciphertext) < aes.BlockSize {
		log.Panicln("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad the decrypted plaintext
	plaintext := unpad(ciphertext)

	return string(plaintext)
}

// mdHash generates an MD5 hash of the input text.
func mdHash(plainText string) []byte {
	hash := md5.Sum([]byte(plainText))
	return hash[:]
}

// pad pads the input to a multiple of the block size.
func pad(input []byte, blockSize int) []byte {
	padding := blockSize - len(input)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padText...)
}

// unpad removes the padding from the input.
func unpad(input []byte) []byte {
	length := len(input)
	unpadding := int(input[length-1])
	return input[:(length - unpadding)]
}
