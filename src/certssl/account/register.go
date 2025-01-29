package account

import (
	"fmt"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

func register(client *lego.Client) (*registration.Resource, error) {
	regOption := registration.RegisterOptions{
		TermsOfServiceAgreed: true,
	}

	reg, err := client.Registration.Register(regOption)
	if err != nil {
		return nil, fmt.Errorf("register failed: %s", err.Error())
	}

	return reg, nil
}
