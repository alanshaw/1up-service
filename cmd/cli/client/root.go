package client

import (
	"github.com/alanshaw/1up-service/cmd/cli/client/debug"
	"github.com/alanshaw/1up-service/cmd/cli/client/provider"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/cobra"
)

var log = logging.Logger("cmd/client")

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "Interact with the 1up service via UCAN invocations",
}

func init() {
	Cmd.AddCommand(debug.Cmd)
	Cmd.AddCommand(provider.Cmd)
}
