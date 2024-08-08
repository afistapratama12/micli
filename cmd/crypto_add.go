package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var cryptoAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "add pair to crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("under maintenance")
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoAddCmd)
}
