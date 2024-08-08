package cmd

import (
	"fmt"
	"strings"

	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

var cryptoAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add pair to crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.PrintErr("pair is required at least 1")
			return
		}

		cryptoView := src.NewCrypto()

		err := cryptoView.AddNewPair(args)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		fmt.Println("add pair successfully")
		fmt.Printf("new pair added: %s\n", strings.Join(args, ", "))
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoAddCmd)
}
