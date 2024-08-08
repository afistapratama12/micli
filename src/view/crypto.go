package view

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/afistapratama12/micli/constants"
	"github.com/afistapratama12/micli/src/model"
	"github.com/afistapratama12/micli/src/repo"
	"github.com/afistapratama12/micli/src/service"
	"github.com/afistapratama12/micli/src/utils"
	"github.com/jedib0t/go-pretty/v6/table"
)

type CryptoView struct {
	CryptoService service.ICrypto
}

func NewCryptoView(CryptoService service.ICrypto) *CryptoView {
	return &CryptoView{
		CryptoService: CryptoService,
	}
}

type PairArg struct {
	Args   string
	Exist  bool
	ErrMsg string
}

func (v *CryptoView) AddNewPair(pairs []string) error {
	// check for valid pair using map
	// if exist value change to true
	var mapPairs = make(map[string]PairArg)

	for idx, pair := range pairs {
		if !strings.Contains(pair, "/") && !strings.Contains(pair, "_") {
			return fmt.Errorf("pair is invalid in argument no %d, format must be have '/' or '_' for example: BTC/USD, BTC_USD or btc/usdt, btc_usdt ...", idx+1)
		}

		key := strings.Replace(pair, "/", "_", -1)

		mapPairs[key] = PairArg{
			Args:  pair,
			Exist: false,
		}
	}

	listSymbol, err := v.CryptoService.GetAllSymbol()
	if err != nil {
		return err
	}

	for _, data := range listSymbol.Data.List {
		if _, ok := mapPairs[data.Symbol]; ok {
			if data.Type > 1 {
				mapPairs[data.Symbol] = PairArg{
					Args:   mapPairs[data.Symbol].Args,
					Exist:  false,
					ErrMsg: "currently not support for market with data in other binance",
				}
				continue
			}

			mapPairs[data.Symbol] = PairArg{
				Args:  mapPairs[data.Symbol].Args,
				Exist: true,
			}
		}
	}

	var validPairs []string

	for key, val := range mapPairs {
		if len(val.ErrMsg) > 0 {
			return fmt.Errorf("error pair %s: %s", val.Args, val.ErrMsg)
		}

		if !val.Exist {
			return fmt.Errorf("error pair %s: not found in market list", val.Args)
		}

		validPairs = append(validPairs, key)
	}

	f, cache, err := utils.ReadFileCache(constants.CACHE_FILE)
	if err != nil {
		return err
	}

	defer f.Close()

	cache.ListPair = append(cache.ListPair, validPairs...)

	err = utils.ModifyFileCache(f, []byte(strings.Join(cache.ListPair, ";")))
	if err != nil {
		return err
	}

	return nil
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

// TODO: this code still not stable
// need to refactor in temp variable every data change
func (v *CryptoView) GetLiveCryptoMarket() error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM) // Notify the channel when an interrupt or SIGTERM signal is received.
	done := make(chan bool, 1)                         // Create a done channel to signal the goroutine to stop.

	_, cache, err := utils.ReadFileCache(constants.CACHE_FILE)
	if err != nil {
		return err
	}

	listPair := cache.ListPair

	var mapPairWS = utils.CreateMapPairWS(listPair)
	dataMap, err := v.CryptoService.CallFirst(listPair)
	if err != nil {
		return err
	}

	// create view first
	tableRows := make([]table.Row, 0)
	for idx, pair := range listPair {
		data := dataMap[pair]
		tableRows = append(tableRows, WriteRow(idx, data, idx < len(listPair)-1)...)
	}

	NewTableMarket(tableRows)

	var wsBase = "wss://stream-cloud.binanceru.net/ws/"
	var wsParam []string
	for _, pair := range listPair {
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
				tableRows := make([]table.Row, 0)
				for idx, pair := range listPair {
					var data model.Result
					mu.Lock()
					data = dataMap[pair]
					mu.Unlock()

					tableRows = append(tableRows, WriteRow(idx, data, idx < len(listPair)-1)...)
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

func WriteRow(idx int, data model.Result, addLine bool) []table.Row {
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
		idx + 1,
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
