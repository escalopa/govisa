package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Encrypter struct {
	secret string
}

func NewEncrypter(secret string) (*Encrypter, error) {
	if len(secret) != 32 {
		return nil, aes.KeySizeError(len(secret))
	}
	return &Encrypter{secret: secret}, nil
}

func (e *Encrypter) Encrypt(text string) (string, error) {
	c, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return encode(ciphertext), nil
}

func (e *Encrypter) Decrypt(text string) (string, error) {
	ciphertext, err := decode(text)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
