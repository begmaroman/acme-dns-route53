package acmstore

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

// retrieveServerCertificate retrieves the server certificate from the given PEM encoded list
func retrieveServerCertificate(list []byte) ([]byte, error) {
	var blocks []*pem.Block
	for {
		var certDERBlock *pem.Block
		certDERBlock, list = pem.Decode(list)
		if certDERBlock == nil {
			break
		}

		if certDERBlock.Type == "CERTIFICATE" {
			blocks = append(blocks, certDERBlock)
		}
	}

	crt := bytes.NewBuffer(nil)
	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse certificate")
		}

		if !cert.IsCA {
			pem.Encode(crt, block)
			break
		}
	}

	return crt.Bytes(), nil
}
