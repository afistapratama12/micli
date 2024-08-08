package cmd

import (
	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoCmd = &cobra.Command{
	Use:     "crypto",
	Aliases: []string{"c"},
	Short:   "Get crypto market information",
	Run: func(cmd *cobra.Command, args []string) {
		cryptoView := src.NewCrypto()

		err := cryptoView.GetLiveCryptoMarket()
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cryptoCmd)
}
