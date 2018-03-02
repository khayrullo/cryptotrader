package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/crankykernel/ctrader/cmd/common"
	"gitlab.com/crankykernel/ctrader/quadriga"
	"github.com/spf13/viper"
)

var quadrigaGetCmd = &cobra.Command{
	Use: "get",
	Run: func(cmd *cobra.Command, args []string) {
		common.Get(quadriga.NewClient(
			viper.GetString("quadriga.api.client-id"),
			viper.GetString("quadriga.api.key"),
			viper.GetString("quadriga.api.secret")), args)
	},
}

func init() {
	quadrigaCmd.AddCommand(quadrigaGetCmd)
}
