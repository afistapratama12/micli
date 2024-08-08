package view

import (
	"fmt"
	"os"

	"github.com/afistapratama12/micli/src/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

// base view code

func NewTableMarket(tableRows []table.Row) {
	utils.RunCmd("clear")

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"#", "pair", "Current Price", "high bid / buy", "low ask / sell", "last buy/sell volume (USD)", "Last Update"})

	t.AppendRows(tableRows)
	t.AppendSeparator()
	t.Render()

	fmt.Printf("\nPress Ctrl+C to exit...\n")
}
