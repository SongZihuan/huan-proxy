package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ReadPrivateKey(data []byte) (crypto.PrivateKey, error) {
	// 解析PEM块
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// 根据PEM块类型解析私钥
	var err error
	var privateKey crypto.PrivateKey
	if block.Type == "RSA PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "EC PRIVATE KEY" {
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	} else {
		return nil, fmt.Errorf("unknown private key type: %s", block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %s", err.Error())
	} else if privateKey == nil {
		return nil, fmt.Errorf("failed to parse private ket: return nil, unknown reason")
	}

	return privateKey, nil
}

func EncodePrivateKeyToPEM(privateKey crypto.PrivateKey) ([]byte, error) {
	var pemType string
	var bytes []byte
	var err error

	switch priv := privateKey.(type) {
	case *rsa.PrivateKey:
		pemType = "RSA PRIVATE KEY"
		bytes = x509.MarshalPKCS1PrivateKey(priv)
	case *ecdsa.PrivateKey:
		pemType = "EC PRIVATE KEY"
		bytes, err = x509.MarshalECPrivateKey(priv)
	default:
		return nil, fmt.Errorf("unsupported private key type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %s", err.Error())
	}

	block := &pem.Block{
		Type:  pemType,
		Bytes: bytes,
	}

	return pem.EncodeToMemory(block), nil
}
