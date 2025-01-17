package applycert

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl/filename"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"os"
	"path"
)

func ReadLocalCertificateAndPrivateKey(basedir string, domain string) (crypto.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	dir := path.Join(basedir, domain)
	cert, err := readCertificate(dir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read certificate failed: %s", err.Error())
	}

	cacert, err := readCACertificate(dir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read certificate failed: %s", err.Error())
	}

	privateKey, err := readPrivateKey(dir)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	return privateKey, cert, cacert, nil
}

func readCertificate(dir string) (*x509.Certificate, error) {
	filepath := path.Join(dir, filename.FileCertificate)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}

	cert, err := utils.ReadCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser certificate file: %v", err)
	}

	return cert, nil
}

func readCACertificate(dir string) (*x509.Certificate, error) {
	filepath := path.Join(dir, filename.FileIssuerCertificate)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}

	cert, err := utils.ReadCertificate(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser certificate file: %v", err)
	}

	return cert, nil
}

func readPrivateKey(dir string) (crypto.PrivateKey, error) {
	filepath := path.Join(dir, filename.FilePrivateKey)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}

	privateKey, err := utils.ReadPrivateKey(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parser key file: %v", err)
	}

	return privateKey, nil
}
