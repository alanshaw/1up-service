package provider

import (
	"github.com/alanshaw/1up-service/cmd/cli/client/provider/weight"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage storage providers",
}

func init() {
	Cmd.AddCommand(deregisterCmd)
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(registerCmd)
	Cmd.AddCommand(weight.Cmd)
}
