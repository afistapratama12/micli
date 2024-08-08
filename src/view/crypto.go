package view

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/afistapratama12/micli/src/model"
	"github.com/afistapratama12/micli/src/repo"
	"github.com/afistapratama12/micli/src/service"
	"github.com/afistapratama12/micli/src/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

var listPairMarket = []string{
	"BTC_USDT", "ETH_USDT", "BNB_USDT", "SOL_USDT", "XRP_USDT", "DOGE_USDT", "ADA_USDT", "TRX_USDT", "SHIB_USDT", "DOT_USDT",
}

var mapPairWS = map[string]string{
	"BTCUSDT":  "BTC_USDT",
	"ETHUSDT":  "ETH_USDT",
	"BNBUSDT":  "BNB_USDT",
	"SOLUSDT":  "SOL_USDT",
	"XRPUSDT":  "XRP_USDT",
	"DOGEUSDT": "DOGE_USDT",
	"ADAUSDT":  "ADA_USDT",
	"TRXUSDT":  "TRX_USDT",
	"SHIBUSDT": "SHIB_USDT",
	"DOTUSDT":  "DOT_USDT",
}

type CryptoView struct {
	CryptoService service.ICrypto
}

func NewCryptoView(CryptoService service.ICrypto) *CryptoView {
	return &CryptoView{
		CryptoService: CryptoService,
	}
}

func (v *CryptoView) GetAllListPair() error {
	listPair, err := v.CryptoService.GetAllListPair()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "pair USD", "pair BTC", "pair BNB", "pair IDR", "pair other"})

	tableRows := make([]table.Row, 0)
	for idx, pair := range listPair {
		pairOther := ""

		if len(pair.Others) > 0 {
			pairOther = strings.Join(pair.Others, ", ")
		}

		tableRows = append(tableRows, table.Row{
			idx + 1,
			pair.USD,
			pair.BTC,
			pair.BNB,
			pair.IDR,
			pairOther,
		})
	}

	t.AppendRows(tableRows)
	t.AppendSeparator()
	t.Render()

	return nil
}

func (v *CryptoView) GetLiveCryptoMarket() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Notify the channel when an interrupt or SIGTERM signal is received.
	done := make(chan bool, 1)                         // Create a done channel to signal the goroutine to stop.

	dataMap, err := v.CryptoService.CallFirst(listPairMarket)
	if err != nil {
		return err
	}

	// create view first
	tableRows := make([]table.Row, 0)
	for idx, pair := range listPairMarket {
		data := dataMap[pair]
		tableRows = append(tableRows, WriteRow(data, idx < len(listPairMarket)-1)...)
	}

	NewTableMarket(tableRows)

	var wsBase = "wss://stream-cloud.binanceru.net/ws/"
	var wsParam []string
	for _, pair := range listPairMarket {
		p := strings.ToLower(strings.Replace(pair, "_", "", 1))
		wsParam = append(wsParam, []string{p + "@depth", p + "@aggTrade"}...)
	}

	c, err := repo.Stream(wsBase + strings.Join(wsParam, "/"))
	if err != nil {
		return err
	}

	var mu sync.Mutex

	defer c.Close()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				msgType, msg, err := c.ReadMessage()
				if err != nil || msgType == -1 {
					continue
				}

				var res model.Result

				if strings.Contains(string(msg), "aggTrade") {
					var data model.WSAggrTradeData
					err = utils.ReadMessage(msg, &data)
					if err != nil {
						log.Printf("error read message, err = %v\n", err)
						continue
					}

					symbol := mapPairWS[data.Symbol]

					mu.Lock()
					if _, ok := dataMap[symbol]; ok {
						res = dataMap[symbol]
					}
					mu.Unlock()

					res.Pair = symbol
					res.AggrData = data.ToAggrTradeData()

					mu.Lock()
					dataMap[symbol] = res
					mu.Unlock()
				} else if strings.Contains(string(msg), "depthUpdate") {
					var data model.WSDepthData
					err = utils.ReadMessage(msg, &data)
					if err != nil {
						log.Printf("error read message, err = %v\n", err)
						continue
					}

					symbol := mapPairWS[data.Symbol]

					mu.Lock()
					if _, ok := dataMap[symbol]; ok {
						res = dataMap[symbol]
					}
					mu.Unlock()

					res.Pair = symbol
					tempRes := res.DepthData

					res.DepthData = data.ToDepthData()

					if res.DepthData.Bids == nil || len(res.DepthData.Asks) < 1 {
						res.DepthData.Bids = tempRes.Bids
					}

					if res.DepthData.Asks == nil || len(res.DepthData.Asks) < 1 {
						res.DepthData.Asks = tempRes.Asks
					}

					mu.Lock()
					dataMap[symbol] = res
					mu.Unlock()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// log.Println("call view in goroutine get data from ws")

				tableRows := make([]table.Row, 0)
				for idx, pair := range listPairMarket {
					var data model.Result
					mu.Lock()
					data = dataMap[pair]
					mu.Unlock()

					tableRows = append(tableRows, WriteRow(data, idx < len(listPairMarket)-1)...)
				}

				NewTableMarket(tableRows)
				time.Sleep(1 * time.Second)
			}
		}
	}()

	<-sigs       // Block until a signal is received.
	done <- true // Send a signal to the goroutine to stop.

	utils.RunCmd("clear")

	return nil
}

func WriteRow(data model.Result, addLine bool) []table.Row {
	var dateTrade time.Time

	if data.AggrData.Timestamp > 0 {
		dateTrade = time.Unix(data.AggrData.Timestamp/1000, 0)
	}

	var bid, ask string
	if len(data.DepthData.Bids) > 0 {
		bid = data.DepthData.Bids[0][0]
	}
	if len(data.DepthData.Asks) > 0 {
		ask = data.DepthData.Asks[0][0]
	}

	var tableRows = make([]table.Row, 0)

	var (
		aggrPrice, aggrVol float64
	)

	if data.AggrData.Price != "" {
		aggrPrice, _ = strconv.ParseFloat(data.AggrData.Price, 64)
	}

	if data.AggrData.Qty != "" {
		aggrVol, _ = strconv.ParseFloat(data.AggrData.Qty, 64)
	}

	lastTradeVol := fmt.Sprintf("%f", aggrPrice*aggrVol)

	tableRows = append(tableRows, table.Row{
		data.Pair,
		data.AggrData.Price,
		bid,
		ask,
		lastTradeVol,
		dateTrade.Format("2006-01-02 15:04:05"),
	})

	if addLine {
		tableRows = append(tableRows, table.Row{
			"", "", "", "", "",
		})
	}

	return tableRows
}
