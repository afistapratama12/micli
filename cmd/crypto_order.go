package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var cryptoOrderCmd = &cobra.Command{
	Use:     "order",
	Aliases: []string{"o"},
	Short:   "reorder pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("under maintenance")
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoOrderCmd)
}
