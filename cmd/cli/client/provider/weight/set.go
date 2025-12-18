package weight

import (
	"strconv"

	"github.com/alanshaw/1up-service/cmd/cli/client/lib"
	"github.com/alanshaw/libracha/capabilities/provider/weight"
	"github.com/alanshaw/ucantone/did"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set <node-did> <num>",
	Short: "Set the weight of a storage provider",
	Args:  cobra.ExactArgs(2),
	RunE:  doSet,
}

func doSet(cmd *cobra.Command, args []string) error {
	signer, client, err := lib.InitClient()
	cobra.CheckErr(err)

	id, err := did.Parse(args[0])
	cobra.CheckErr(err)

	w, err := strconv.ParseUint(args[1], 10, 64)
	cobra.CheckErr(err)

	inv, err := weight.Set.Invoke(
		signer,
		signer,
		&weight.SetArguments{Provider: id, Weight: w},
		invocation.WithAudience(signer),
	)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Result(),
		func(o ipld.Any) {
			args := weight.SetOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)
			cmd.Println("Provider weight set successfully")
		},
		func(x ipld.Any) {
			cmd.Printf("Invocation failed: %+v\n", x)
		},
	)
	return nil
}
