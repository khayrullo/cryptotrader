package cmd

import (
	"github.com/spf13/cobra"
)

var quadrigaCmd = &cobra.Command{
	Use: "quadriga",
}

func init() {
	rootCmd.AddCommand(quadrigaCmd)
}
