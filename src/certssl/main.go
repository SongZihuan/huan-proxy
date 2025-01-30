package certssl

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl/applycert"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"time"
)

const CertDefaultNewApplyTime = 5 * 24 * time.Hour

func GetCertificateAndPrivateKey(basedir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string) (crypto.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	if email == "" {
		return nil, nil, nil, fmt.Errorf("email is empty")
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, nil, fmt.Errorf("not a valid domain")
	}

	privateKey, cert, cacert, err := applycert.ReadLocalCertificateAndPrivateKey(basedir, domain)
	if err == nil && utils.CheckCertWithDomain(cert, domain) && utils.CheckCertWithTime(cert, 5*24*time.Hour) {
		return privateKey, cert, cacert, nil
	}

	resource, err := applycert.ApplyCert(basedir, email, aliyunAccessKey, aliyunAccessSecret, domain)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("apply cert failed: %s", err.Error())
	} else if resource == nil {
		return nil, nil, nil, fmt.Errorf("read cert failed: private key or certificate (resource) is nil, unknown reason")
	}

	privateKey, err = utils.ReadPrivateKey(resource.PrivateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	cert, err = utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	cacert, err = utils.ReadCertificate(resource.IssuerCertificate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	return privateKey, cert, cacert, nil
}

type NewCert struct {
	PrivateKey        crypto.PrivateKey
	Certificate       *x509.Certificate
	IssuerCertificate *x509.Certificate
	Error             error
}

func WatchCertificate(dir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string, oldCert *x509.Certificate, stopchan chan bool, newchan chan NewCert) error {
	for {
		select {
		case <-stopchan:
			newchan <- NewCert{
				PrivateKey:  nil,
				Certificate: nil,
				Error:       nil,
			}
			close(stopchan)
			return nil
		default:
			privateKey, cert, cacert, err := watchCertificate(dir, email, aliyunAccessKey, aliyunAccessSecret, domain, oldCert)
			if err != nil {
				newchan <- NewCert{
					Error: fmt.Errorf("watch cert failed: %s", err.Error()),
				}
			} else if privateKey != nil && cert != nil && cacert != nil {
				oldCert = cert
				newchan <- NewCert{
					PrivateKey:        privateKey,
					Certificate:       cert,
					IssuerCertificate: cacert,
				}
			}
		}
	}
}

func watchCertificate(dir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string, oldCert *x509.Certificate) (crypto.PrivateKey, *x509.Certificate, *x509.Certificate, error) {
	if email == "" {
		return nil, nil, nil, fmt.Errorf("email is empty")
	}

	if !utils.IsValidEmail(email) {
		return nil, nil, nil, fmt.Errorf("not a valid email")
	}

	if !utils.IsValidDomain(domain) {
		return nil, nil, nil, fmt.Errorf("not a valid domain")
	}

	if utils.CheckCertWithDomain(oldCert, domain) && utils.CheckCertWithTime(oldCert, CertDefaultNewApplyTime) {
		return nil, nil, nil, nil
	}

	resource, err := applycert.ApplyCert(dir, email, aliyunAccessKey, aliyunAccessSecret, domain)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("apply cert fail: %s", err.Error())
	}

	privateKey, err := utils.ReadPrivateKey(resource.PrivateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read private key failed: %s", err.Error())
	}

	cert, err := utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	cacert, err := utils.ReadCertificate(resource.IssuerCertificate)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read cert failed: %s", err.Error())
	}

	return privateKey, cert, cacert, nil
}
