package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var cryptoRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "r"},
	Short:   "remove pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("under maintenance")
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoRemoveCmd)
}
