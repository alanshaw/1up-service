package debug

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug the service",
}

func init() {
	Cmd.AddCommand(echoCmd)
}
