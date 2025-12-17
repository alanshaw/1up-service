package debug

import (
	"fmt"
	"net/url"

	"github.com/alanshaw/1up-service/pkg/capabilities/debug"
	"github.com/alanshaw/1up-service/pkg/config"
	"github.com/alanshaw/ucantone/client"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Echo a message",
	Args:  cobra.ExactArgs(1),
	RunE:  doEcho,
}

func init() {
	//
}

func doEcho(cmd *cobra.Command, args []string) error {
	userCfg, err := config.Load[config.AppConfig]()
	cobra.CheckErr(err)

	appCfg, err := userCfg.ToAppConfig()
	cobra.CheckErr(err)

	inv, err := debug.Echo.Invoke(
		appCfg.Identity.Signer,
		appCfg.Identity.Signer,
		&debug.EchoArguments{
			Message: args[0],
		},
		invocation.WithAudience(appCfg.Identity.Signer),
	)
	cobra.CheckErr(err)

	url, err := url.Parse(fmt.Sprintf("http://%s:%d", appCfg.Server.Host, appCfg.Server.Port))
	cobra.CheckErr(err)

	client, err := client.NewHTTP(url)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Result(),
		func(o ipld.Any) {
			args := debug.EchoOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)
			fmt.Printf("Echo response: %+v\n", args)
		},
		func(x ipld.Any) {
			fmt.Printf("Invocation failed: %v\n", x)
		},
	)

	return nil
}
