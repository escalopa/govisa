package security

import (
	"math/rand"
	"testing"
	"time"
)

func TestEnctyptor(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		text   string
		err    error
	}{{
		name:   "Test encrypt and decrypt",
		secret: randomString(32),
		text:   "Hello World",
		err:    nil,
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewEncrypter(tt.secret)
			if err != nil {
				t.Errorf("NewEncrypter() error = %v", err)
				return
			}
			encrypted, err := e.Encrypt(tt.text)
			if err != nil {
				t.Errorf("Encrypt() error = %v", err)
				return
			}
			decrypted, err := e.Decrypt(encrypted)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}
			if decrypted != tt.text && err == nil {
				t.Errorf("Decrypt() = %v, want %v", decrypted, tt.text)
			}
		})
	}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return string(b)
}
