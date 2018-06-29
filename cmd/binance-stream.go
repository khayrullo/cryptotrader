package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/crankykernel/cryptotrader/binance"
	"log"
	"fmt"
	"encoding/json"
)

var binanceStreamSingle bool

var binanceStreamCmd = &cobra.Command{
	Use:   "stream <stream0> <stream1> ...",
	Short: "Print one or more streams",
	Long: `Connects to the Binance websocket and prints the output of one or more stream
names provided on the command line.
`,
	Run: func(cmd *cobra.Command, args []string) {
		client := binance.NewStreamClient()
		if binanceStreamSingle {
			if err := client.ConnectSingle(args[0]); err != nil {
				log.Fatal("error: ", err)
			}
		} else {
			if err := client.Connect(args...); err != nil {
				log.Fatal("error: ", err)
			}
		}
		log.Println("Connected!")
		for {
			msg, err := client.NextJSON()
			if err != nil {
				log.Fatal("error: ", err)
			}
			txt, err := json.Marshal(msg)
			if err != nil {
				log.Fatal("error: ", err)
			}
			fmt.Printf("%s\n", txt)
		}
	},
}

func init() {
	binanceCmd.AddCommand(binanceStreamCmd)

	flags := binanceStreamCmd.Flags()
	flags.BoolVarP(&binanceStreamSingle, "single", "s", false,
		"Use the single stream endpoint.")
}
