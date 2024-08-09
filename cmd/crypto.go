package cmd

import (
	"strconv"

	"github.com/afistapratama12/micli/src"
	"github.com/spf13/cobra"
)

//TODO: ignore comment below, this for note development

var realtime bool
var interval string

// var isLongDepth bool

var cryptoCmd = &cobra.Command{
	Use:     "crypto",
	Aliases: []string{"c"},
	Short:   "Get crypto market information",
	Run: func(cmd *cobra.Command, args []string) {
		var intervalTime = -1
		var err error

		if interval != "" || len(interval) > 0 {
			if len(interval) < 2 {
				cmd.PrintErr("Interval format must be in number and time unit, example: 1s, 3s, 10s")
			}

			if interval[len(interval)-1] != 's' {
				cmd.PrintErr("only avaible time unit is second, example: 1s, 3s, 10s")
			}

			intervalTime, err = strconv.Atoi(interval[:len(interval)-1])
			if err != nil {
				cmd.PrintErr("Interval format must be in number and time unit, example: 1s, 3s, 10s")
			}
		}

		cryptoView := src.NewCrypto()
		err = cryptoView.GetLiveCryptoMarket(realtime, intervalTime)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cryptoCmd)

	cryptoCmd.Flags().BoolVarP(&realtime, "realtime", "r", false, "Realtime crypto market information")
	cryptoCmd.Flags().StringVarP(&interval, "interval", "i", "", "Interval refresh data crypto market information, default 1s")
	cryptoCmd.Flags().Lookup("realtime").NoOptDefVal = "true"
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
