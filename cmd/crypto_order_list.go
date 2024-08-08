package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var cryptoOrderListCmd = &cobra.Command{
	Use:     "order-list",
	Aliases: []string{"ol"},
	Short:   "list order pair in crypto market view",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("under maintenance")
	},
}

func init() {
	cryptoCmd.AddCommand(cryptoOrderListCmd)
}
