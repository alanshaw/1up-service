package provider

import (
	"github.com/alanshaw/1up-service/pkg/config"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a storage node with the service",
	Args:  cobra.ExactArgs(2),
	RunE:  doRegister,
}

func init() {
	//
}

func doRegister(cmd *cobra.Command, args []string) error {
	userCfg, err := config.Load[config.AppConfig]()
	cobra.CheckErr(err)

	_, err = userCfg.ToAppConfig()
	cobra.CheckErr(err)

	// appCfg.Identity.Signer
	return nil
}
