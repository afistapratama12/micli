package cmd

import (
	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

//TODO: ignore comment below, this for note development

// var realtime bool
// var isLongDepth bool

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

		// log.Println("realtime: ", realtime)
		// realtime, _ := cmd.Flags().GetBool("realtime")
	},
}

func init() {
	// cryptoCmd.Flags().BoolVar(&realtime, "realtime", false, "Get realtime crypto market information")
	rootCmd.AddCommand(cryptoCmd)

	// cryptoCmd.Flags().BoolVar(&isLongDepth, "long-depth", false, "Get 3 high bid and 3 low asks info") // 3 depth bid and 3 depth ask

	// cryptoCmd.Flags().BoolVar(false, "rt", "", "config for realtime crypto market")
}

// command

// micli crypto --realtime / -rt
// micli crypto --long-depth / -ld

// micli crypto [done, option still not implement]
// micli crypto list [done]
// micli crypto add [pair] [pair2] [pair...] --order-at=1 [done, option still not implement]
// micli crypto order [pair] [order-number]
// micli crypto order-list
// micli crypto remove [pair] [pair2] [pair...]

// note pair can be lower or upper
// but must have symbol "/" or "_" in between

// --order = must be in number and cannot less than 1, and more than list existing market crypto

//
