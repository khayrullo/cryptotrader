// The MIT License (MIT)
//
// Copyright (c) 2018 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var kucoinCmd = &cobra.Command{
	Use: "kucoin",
	Long: `KuCoin API Tools

Configuration File Parameters:
  - kucoin.api.key
  - kucoin.api.secret

Environment Variables:
  - KUCOIN_API_KEY
  - KUCOIN_API_SECRET`,
}

func init() {
	kucoinCmd.PersistentFlags().String("api-key", "", "KuCoin API key")
	viper.BindPFlag("kucoin.api.key", kucoinCmd.PersistentFlags().Lookup("api-key"))
	viper.BindEnv("kucoin.api.key", "KUCOIN_API_KEY")

	kucoinCmd.PersistentFlags().String("api-secret", "", "KuCoin API secret")
	viper.BindPFlag("kucoin.api.secret", kucoinCmd.PersistentFlags().Lookup("api-secret"))
	viper.BindEnv("kucoin.api.secret", "KUCOIN_API_SECRET")

	rootCmd.AddCommand(kucoinCmd)
}
