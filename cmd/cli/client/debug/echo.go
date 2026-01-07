package debug

import (
	"github.com/alanshaw/1up-service/cmd/cli/client/lib"
	"github.com/alanshaw/libracha/capabilities/debug"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo <message>",
	Short: "Echo a message",
	Args:  cobra.ExactArgs(1),
	RunE:  doEcho,
}

func doEcho(cmd *cobra.Command, args []string) error {
	signer, client, err := lib.InitClient()
	cobra.CheckErr(err)

	inv, err := debug.Echo.Invoke(
		signer,
		signer,
		&debug.EchoArguments{
			Message: args[0],
		},
	)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Out(),
		func(o ipld.Any) {
			args := debug.EchoOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)
			cmd.Printf("Echo response: %+v\n", args)
		},
		func(x ipld.Any) {
			cmd.Printf("Invocation failed: %+v\n", x)
		},
	)
	return nil
}
