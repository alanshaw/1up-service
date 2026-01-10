package provider

import (
	"fmt"

	"github.com/alanshaw/1up-service/cmd/cli/client/lib"
	"github.com/alanshaw/libracha/capabilities/provider"
	"github.com/alanshaw/ucantone/execution"
	"github.com/alanshaw/ucantone/ipld"
	"github.com/alanshaw/ucantone/ipld/datamodel"
	"github.com/alanshaw/ucantone/result"
	"github.com/alanshaw/ucantone/ucan/invocation"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List registered storage nodes",
	Args:    cobra.NoArgs,
	RunE:    doList,
}

func doList(cmd *cobra.Command, args []string) error {
	signer, client, err := lib.InitClient()
	cobra.CheckErr(err)

	inv, err := provider.List.Invoke(
		signer,
		signer,
		&provider.ListArguments{},
		invocation.WithAudience(signer),
	)
	cobra.CheckErr(err)

	res, err := client.Execute(execution.NewRequest(cmd.Context(), inv))
	cobra.CheckErr(err)

	result.MatchResultR0(
		res.Receipt().Out(),
		func(o ipld.Any) {
			args := provider.ListOK{}
			err := datamodel.Rebind(datamodel.NewAny(o), &args)
			cobra.CheckErr(err)

			if len(args.Providers) == 0 {
				cmd.Println("No providers registered")
				return
			}

			table := lib.NewTable(cmd.OutOrStdout())
			table.SetHeader([]string{"ID", "Weight", "URL"})
			for _, p := range args.Providers {
				table.Append([]string{p.Provider.String(), fmt.Sprintf("%d", p.Weight), p.Endpoint})
			}
			table.Render()
		},
		func(x ipld.Any) {
			cmd.Printf("Invocation failed: %+v\n", x)
		},
	)
	return nil
}
