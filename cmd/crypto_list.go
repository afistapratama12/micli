package cmd

import (
	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Get list of crypto market pair",
	Run: func(cmd *cobra.Command, args []string) {
		cryptoView := src.NewCrypto()
		err := cryptoView.GetAllListPair()
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoListCmd)
}
