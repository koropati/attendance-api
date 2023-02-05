package cryptos

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

type Cryptos interface {
	Encrypt(inputString string) (output string)
	Decrypt(inputString string) (output string)
	EncryptAES256(inputString string) (output string)
	DecryptAES256(inputString string) (output string)
}

type cryptos struct {
	c cipher.Stream
	a cipher.AEAD
}

func New(c cipher.Stream, a cipher.AEAD) Cryptos {
	return &cryptos{
		c: c,
		a: a,
	}
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

func (c *cryptos) EncryptAES256(inputString string) (output string) {

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, c.a.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := c.a.Seal(nonce, nonce, []byte(inputString), nil)
	return fmt.Sprintf("%x", ciphertext)
}

func (c *cryptos) DecryptAES256(inputString string) (output string) {

	enc, _ := hex.DecodeString(inputString)
	//Get the nonce size
	nonceSize := c.a.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := c.a.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%x", plaintext)
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
