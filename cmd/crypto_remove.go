package cmd

import (
	"fmt"
	"strings"

	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   "remove pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErr("invalid argument\nformat: micli crypto remove [pair] [pair2] [pair...]")
			return
		}

		cryptoView := src.NewCrypto()
		err := cryptoView.RemovePair(args)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		fmt.Printf("remove pair success\n")
		fmt.Printf("pair removed: %s\n", strings.Join(args, ", "))
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoRemoveCmd)
}
