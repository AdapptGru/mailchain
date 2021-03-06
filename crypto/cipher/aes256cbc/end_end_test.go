package aes256cbc

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/crypto/secp256k1/secp256k1test"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name                string
		recipientPublicKey  crypto.PublicKey
		recipientPrivateKey crypto.PrivateKey
		data                []byte
		wantEncryptErr      bool
		wantDecryptErr      bool
		wantDecrypt         bool
	}{
		{
			"success-to-sofia-short-text",
			secp256k1test.SofiaPublicKey,
			secp256k1test.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			false,
			false,
			true,
		},
		{
			"success-to-sofia-medium-text",
			secp256k1test.SofiaPublicKey,
			secp256k1test.SofiaPrivateKey,
			[]byte("Hi Sofia, this is a little bit of a longer message to make sure there are no problems"),
			false,
			false,
			true,
		},
		{
			"success-to-charlotte-short-text",
			secp256k1test.CharlottePublicKey,
			secp256k1test.CharlottePrivateKey,
			[]byte("Hi Charlotte"),
			false,
			false,
			true,
		},
		{
			"success-to-charlotte-medium-text",
			secp256k1test.CharlottePublicKey,
			secp256k1test.CharlottePrivateKey,
			[]byte("Hi Charlotte, this is a little bit of a longer message to make sure there are no problems"),
			false,
			false,
			true,
		},
		{
			"err-sofia-with-charlotte",
			secp256k1test.SofiaPublicKey,
			secp256k1test.CharlottePrivateKey,
			[]byte("Hi Sofia"),
			false,
			true,
			false,
		},
		{
			"err-charlotte-with-sofia",
			secp256k1test.CharlottePublicKey,
			secp256k1test.SofiaPrivateKey,
			[]byte("Hi Sofia"),
			false,
			true,
			false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			encrypter := Encrypter{rand.Reader, tt.recipientPublicKey}
			encrypted, err := encrypter.Encrypt(tt.data)
			if (err != nil) != tt.wantEncryptErr {
				t.Errorf("Encrypt() error = %v, wantEncryptErr %v", err, tt.wantEncryptErr)
				return
			}
			assert.NotNil(encrypted)

			decrypter := Decrypter{tt.recipientPrivateKey}
			decrypted, err := decrypter.Decrypt(encrypted)
			if (err != nil) != tt.wantDecryptErr {
				t.Errorf("Decrypt() error = %v, wantDecryptErr %v", err, tt.wantDecryptErr)
				return
			}

			if bytes.Equal(tt.data, []byte(decrypted)) != tt.wantDecrypt {
				t.Errorf("Decrypt() result = %v, wantDecrypt %v", err, tt.wantDecrypt)
				return
			}
		})
	}
}
