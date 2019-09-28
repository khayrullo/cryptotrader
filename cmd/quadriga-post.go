package cmd

import (
	"github.com/spf13/cobra"
	"cryptotrader/cmd/common"
	"cryptotrader/quadriga"
	"github.com/spf13/viper"
)

var quadrigaPostCmd = &cobra.Command{
	Use: "post",
	Run: func(cmd *cobra.Command, args []string) {
		common.Post(quadriga.NewClient(
			viper.GetString("quadriga.api.client-id"),
			viper.GetString("quadriga.api.key"),
			viper.GetString("quadriga.api.secret")), args)
	},
}

func init() {
	quadrigaCmd.AddCommand(quadrigaPostCmd)
}
