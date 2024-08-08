package service

import (
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/afistapratama12/micli/src/model"
	"github.com/afistapratama12/micli/src/repo"
)

type ICrypto interface {
	GetAllListPair() (res []model.ViewPairMarket, err error)
	CallFirst(listPairMarket []string) (map[string]model.Result, error)
}

type crypto struct {
	CryptoRepo repo.ICrypto
}

func NewCryptoService(cryptoRepo repo.ICrypto) ICrypto {
	return &crypto{
		CryptoRepo: cryptoRepo,
	}
}

func (s *crypto) GetAllListPair() (res []model.ViewPairMarket, err error) {
	listData, err := s.CryptoRepo.GetAllListPair()
	if err != nil {
		return
	}

	var mappingPair = make(map[string]model.ViewPairMarket)

	for _, data := range listData.Data.List {
		symb := strings.Split(data.Symbol, "_")

		var pair model.ViewPairMarket

		if _, ok := mappingPair[symb[0]]; !ok {
			mappingPair[symb[0]] = model.ViewPairMarket{}
		}

		pair = mappingPair[symb[0]]

		switch {
		case ListContains(data.Symbol, "_USDT", "_TUSD", "_USDC"):
			pair.USD = data.Symbol
		case ListContains(data.Symbol, "_IDR"):
			pair.IDR = data.Symbol
		case strings.Contains(data.Symbol, "_BTC"):
			pair.BTC = data.Symbol
		case strings.Contains(data.Symbol, "_BNB"):
			pair.BNB = data.Symbol
		default:
			pair.Others = append(pair.Others, data.Symbol)
		}

		mappingPair[symb[0]] = pair
	}

	for _, pair := range mappingPair {
		res = append(res, pair)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].USD < res[j].USD && (res[i].USD != "" || res[j].USD != "")
	})

	return
}

func ListContains(s string, comp ...string) bool {
	for _, c := range comp {
		if strings.Contains(s, c) {
			return true
		}
	}

	return false
}

func (s *crypto) CallFirst(listPairMarket []string) (map[string]model.Result, error) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)               // Limit to 5 concurrent goroutines
	var mu sync.Mutex                           // Mutex to protect shared slice
	var results = make(map[string]model.Result) // Slice to store the results

	for _, pair := range listPairMarket {
		wg.Add(1)
		go func(pair string) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire a token
			defer func() { <-sem }() // Release the token

			var depthData model.DepthData
			var listAggr []model.AggrTradeData

			depthData, err := s.CryptoRepo.GetDepth(pair)
			if err != nil {
				log.Printf("error call depth data, pair %s, err = %v", pair, err)
				panic(err)
			}

			listAggr, err = s.CryptoRepo.GetAggrTrade(pair)
			if err != nil {
				log.Printf("error call aggr Trades, pair %s, err = %v", pair, err)
				panic(err)
			}

			// Collect the result
			mu.Lock()
			results[pair] = model.Result{
				Pair:      pair,
				DepthData: depthData,
				AggrData:  listAggr[0],
			}

			mu.Unlock()

		}(pair)
	}

	wg.Wait()

	return results, nil
}
