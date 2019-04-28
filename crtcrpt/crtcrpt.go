package crtcrpt

import (
	"crypto"
	"encoding/pem"
	"os"

	"github.com/go-acme/lego/certcrypto"
)

// GeneratePrivateKey generates a new private key based on the given key type and stores it to a new created file
func GeneratePrivateKey(file string, keyType certcrypto.KeyType) (crypto.PrivateKey, error) {
	privateKey, err := certcrypto.GeneratePrivateKey(keyType)
	if err != nil {
		return nil, err
	}

	certOut, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer certOut.Close()

	pemKey := certcrypto.PEMBlock(privateKey)
	err = pem.Encode(certOut, pemKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
