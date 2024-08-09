package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoOrderCmd = &cobra.Command{
	Use:     "order",
	Aliases: []string{"o"},
	Short:   "reorder pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		// max arg 2
		if len(args) < 2 {
			cmd.PrintErr("invalid argument\n format: micli crypto order [pair] [order-number]")
			return
		}

		if !strings.Contains(args[0], "/") && !strings.Contains(args[0], "_") {
			cmd.PrintErr("pair is invalid, format must be have '/' or '_' for example: BTC/USD, BTC_USD or btc/usdt, btc_usdt")
			return
		}

		ord, err := strconv.Atoi(args[1])
		if err != nil {
			cmd.PrintErr("invalid order number")
			return
		}

		if ord < 1 {
			cmd.PrintErr("order number must start from 1")
		}

		cryptoView := src.NewCrypto()
		err = cryptoView.ReorderPair(args[0], ord)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		fmt.Printf("reorder pair success\n")
		fmt.Printf("pair: %s, change to order: %d\n", args[0], ord)
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoOrderCmd)
}
