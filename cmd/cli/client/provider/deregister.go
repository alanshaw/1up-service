package provider

import (
	"github.com/alanshaw/1up-service/cmd/cli/client/lib"
	"github.com/alanshaw/1up-service/pkg/capabilities/provider"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/spf13/cobra"
)

var deregisterCmd = &cobra.Command{
	Use:     "deregister <node-did>",
	Aliases: []string{"remove", "rm"},
	Short:   "Deregister a storage node from the service",
	Args:    cobra.ExactArgs(1),
	RunE:    doDeregister,
}

func doDeregister(cmd *cobra.Command, args []string) error {
	signer, client, err := lib.InitClient()
	cobra.CheckErr(err)

	id, err := did.Parse(args[0])
	cobra.CheckErr(err)

	inv, err := provider.Deregister.Invoke(
		signer,
		signer,
		&provider.DeregisterArguments{
			Provider: id,
		},
		invocation.WithAudience(signer),
	)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Result(),
		func(o ipld.Any) {
			args := provider.DeregisterOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)
			cmd.Println("Provider deregistered successfully")
		},
		func(x ipld.Any) {
			cmd.Printf("Invocation failed: %+v\n", x)
		},
	)
	return nil
}
