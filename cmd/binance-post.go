// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/crankykernel/cryptotrader/binance"
	"github.com/spf13/viper"
	"gitlab.com/crankykernel/cryptotrader/cmd/common"
)

var binancePostCmd = &cobra.Command{
	Use: "post",
	Run: func(cmd *cobra.Command, args []string) {
		var client *binance.RestClient
		auth, _ := cmd.Flags().GetBool("auth")
		if auth {
			client = binance.NewAuthenticatedClient(
				viper.GetString("binance.api.key"),
				viper.GetString("binance.api.secret"))
		} else {
			client = binance.NewAnonymousClient()
		}
		common.Post(client, args)
	},
}

func init() {
	flags := binancePostCmd.Flags()

	flags.Bool("auth", false, "Send authenticated request")

	binanceCmd.AddCommand(binancePostCmd)
}
