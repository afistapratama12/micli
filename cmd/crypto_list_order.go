package cmd

import (
	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoListOrderCmd = &cobra.Command{
	Use:     "list-order",
	Aliases: []string{"lo"},
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
	cryptoCmd.AddCommand(cryptoListOrderCmd)
}
