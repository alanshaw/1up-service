package provider

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage storage providers",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(registerCmd)
}
