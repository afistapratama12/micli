package cmd

import (
	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoOrderListCmd = &cobra.Command{
	Use:     "order-list",
	Aliases: []string{"ol"},
	Short:   "list order pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		cryptoView := src.NewCrypto()
		err := cryptoView.GetOrderList()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoOrderListCmd)
}
