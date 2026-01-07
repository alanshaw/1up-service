package provider

import (
	"net/url"

	"github.com/alanshaw/1up-service/cmd/cli/client/lib"
	"github.com/alanshaw/libracha/capabilities/provider"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:     "register <node-did> <node-url>",
	Aliases: []string{"add"},
	Short:   "Register a storage node with the service",
	Args:    cobra.ExactArgs(2),
	RunE:    doRegister,
}

func doRegister(cmd *cobra.Command, args []string) error {
	signer, client, err := lib.InitClient()
	cobra.CheckErr(err)

	id, err := did.Parse(args[0])
	cobra.CheckErr(err)

	endpoint, err := url.Parse(args[1])
	cobra.CheckErr(err)

	inv, err := provider.Register.Invoke(
		signer,
		signer,
		&provider.RegisterArguments{
			Provider: id,
			Endpoint: endpoint.String(),
		},
		invocation.WithAudience(signer),
	)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Out(),
		func(o ipld.Any) {
			args := provider.RegisterOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)
			cmd.Println("Provider registered successfully")
		},
		func(x ipld.Any) {
			cmd.Printf("Invocation failed: %+v\n", x)
		},
	)
	return nil
}
