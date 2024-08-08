package cmd

import (
	"fmt"

	"github.com/afistapratama12/micli/constants"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of micli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("micli %s\n", constants.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
