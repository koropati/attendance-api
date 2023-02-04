package cryptos

import (
	"crypto/cipher"
	"encoding/base64"
	"log"
)

type Cryptos interface {
	Encrypt(inputString string) (output string)
	Decrypt(inputString string) (output string)
}

type cryptos struct {
	c cipher.Stream
}

func New(c cipher.Stream) Cryptos {
	return &cryptos{c: c}
}

func (c *cryptos) Encrypt(inputString string) (output string) {
	plainText := []byte(inputString)
	cipherText := make([]byte, len(plainText))

	c.c.XORKeyStream(cipherText, plainText)
	output = string(plainText)
	return
}

func (c *cryptos) Decrypt(inputString string) (output string) {
	ciphertext := Decode(inputString)
	plainText := make([]byte, len(ciphertext))

	c.c.XORKeyStream(plainText, ciphertext)
	output = string(plainText)
	return
}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Printf("[Error] [Decode] E: %v", err)
	}
	return data
}
