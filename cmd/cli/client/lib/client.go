package lib

import (
	"fmt"
	"net/url"

	"github.com/alanshaw/1up-service/pkg/config"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/principal"
)

func InitClient() (principal.Signer, execution.Executor, error) {
	userCfg, err := config.Load[config.AppConfig]()
	if err != nil {
		return nil, nil, err
	}

	appCfg, err := userCfg.ToAppConfig()
	if err != nil {
		return nil, nil, err
	}

	url, err := url.Parse(fmt.Sprintf("http://%s:%d", appCfg.Server.Host, appCfg.Server.Port))
	if err != nil {
		return nil, nil, err
	}

	client, err := client.NewHTTP(url)
	if err != nil {
		return nil, nil, err
	}

	return appCfg.Identity.Signer, client, nil
}
