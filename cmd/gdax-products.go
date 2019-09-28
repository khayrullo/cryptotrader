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
	"fmt"

	"github.com/khayrullo/cryptotrader/gdax"
	"github.com/spf13/cobra"
	"log"
)

var gdaxProductsCmd = &cobra.Command{
	Use:   "products",
	Short: "List GDAX products (pairs)",
	Run: func(cmd *cobra.Command, args []string) {
		products, err := gdax.NewApiClient().Products()
		if err != nil {
			log.Fatal("error: ", err)
		}
		for _, product := range products {
			fmt.Printf("ID: %s: Name: %s\n", product.Id, product.DisplayName)
		}
	},
}

func init() {
	log.SetFlags(0)
	gdaxCmd.AddCommand(gdaxProductsCmd)
}
