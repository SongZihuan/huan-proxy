package applycert

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/certssl/account"
	"github.com/SongZihuan/huan-proxy/src/logger"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"time"
)

const DefaultCertTimeout = 30 * 24 * time.Hour
const DefaultCertType = certcrypto.RSA4096

func ApplyCert(basedir string, email string, aliyunAccessKey string, aliyunAccessSecret string, domain string) (*certificate.Resource, error) {
	if domain == "" || !utils.IsValidDomain(domain) {
		return nil, fmt.Errorf("domain is invalid")
	}

	user, err := account.LoadAccount(basedir, email)
	if err != nil {
		logger.Infof("load local account failed, register a new on for %s: %s\n", email, err.Error())
		user, err = account.NewAccount(basedir, email)
		if err != nil {
			return nil, fmt.Errorf("generate new user failed: %s", err.Error())
		}
	}

	config := lego.NewConfig(user)
	config.Certificate.KeyType = DefaultCertType
	config.Certificate.Timeout = DefaultCertTimeout
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("new client failed: %s", err.Error())
	}

	aliyunDnsConfig := alidns.NewDefaultConfig()
	aliyunDnsConfig.APIKey = aliyunAccessKey
	aliyunDnsConfig.SecretKey = aliyunAccessSecret

	provider, err := alidns.NewDNSProviderConfig(aliyunDnsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize AliDNS provider: %s", err.Error())
	}

	err = client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return nil, fmt.Errorf("set challenge dns1 provider failed: %s", err.Error())
	}

	reg, err := user.Register(client)
	if err != nil {
		return nil, fmt.Errorf("get account failed: %s", err.Error())
	} else if reg == nil {
		return nil, fmt.Errorf("get account failed: return nil account.resurce, unknown reason")
	}

	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	resource, err := client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("obtain certificate failed: %s", err.Error())
	}

	err = user.SaveAccount()
	if err != nil {
		return nil, fmt.Errorf("save account error after obtain: %s", err.Error())
	}

	cert, err := utils.ReadCertificate(resource.Certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate: %s", err.Error())
	}

	err = writerWithDate(basedir, cert, resource)
	if err != nil {
		return nil, fmt.Errorf("writer certificate backup failed: %s", err.Error())
	}

	err = writer(basedir, cert, resource)
	if err != nil {
		return nil, fmt.Errorf("writer certificate failed: %s", err.Error())
	}

	return resource, nil
}
