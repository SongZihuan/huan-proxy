package applycert

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl/filename"
	"github.com/go-acme/lego/v4/certificate"
	"os"
	"path"
)

func writerWithDate(basedir string, cert *x509.Certificate, resource *certificate.Resource) error {
	domain := cert.Subject.CommonName
	if domain == "" && len(cert.DNSNames) == 0 {
		return fmt.Errorf("no domains in certificate")
	}
	domain = cert.DNSNames[0]

	year := fmt.Sprintf("%d", cert.NotBefore.Year())
	month := fmt.Sprintf("%d", cert.NotBefore.Month())
	day := fmt.Sprintf("%d", cert.NotBefore.Day())

	backupdir := path.Join(basedir, "cert-backup", domain, year, month, day, cert.NotBefore.Format("2006-01-02-15:04:05"))
	err := os.MkdirAll(backupdir, 0775)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	data, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(backupdir, filename.FileResource), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func writer(basedir string, cert *x509.Certificate, resource *certificate.Resource) error {
	domain := cert.Subject.CommonName
	if domain == "" && len(cert.DNSNames) == 0 {
		return fmt.Errorf("no domains in certificate")
	}
	domain = cert.DNSNames[0]

	dir := path.Join(basedir, domain)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", dir, err.Error())
	}

	err = os.WriteFile(path.Join(dir, filename.FilePrivateKey), resource.PrivateKey, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCertificate), resource.Certificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileIssuerCertificate), resource.IssuerCertificate, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(dir, filename.FileCSR), resource.CSR, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
