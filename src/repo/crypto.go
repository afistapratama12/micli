package repo

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/afistapratama12/micli/src/model"
)

type ICrypto interface {
	GetAllListPair() (*model.MarketView, error)
	GetDepth(pair string) (model.DepthData, error)
	GetAggrTrade(pair string) ([]model.AggrTradeData, error)
}

type crypto struct {
}

func NewCryptoRepo() ICrypto {
	return &crypto{}
}

func (r *crypto) GetAllListPair() (*model.MarketView, error) {
	var res model.MarketView
	err := Call(model.ReqData{
		Method: http.MethodGet,
		Url:    "https://www.tokocrypto.com/open/v1/common/symbols",
	}, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *crypto) GetDepth(pair string) (model.DepthData, error) {
	var depthData model.DepthData

	err := Call(model.ReqData{
		Method: http.MethodGet,
		Url:    "https://api.binance.com/api/v3/depth",
		Params: url.Values{
			"symbol": []string{strings.Replace(pair, "_", "", 1)},
			"limit":  []string{"5"},
		},
	}, &depthData)

	if err != nil {
		return model.DepthData{}, err
	}

	return depthData, nil
}

func (r *crypto) GetAggrTrade(pair string) ([]model.AggrTradeData, error) {
	var listAggr []model.AggrTradeData

	err := Call(model.ReqData{
		Method: http.MethodGet,
		Url:    "https://api.binance.com/api/v3/aggTrades",
		Params: url.Values{
			"symbol": []string{strings.Replace(pair, "_", "", 1)},
			"limit":  []string{"1"},
		},
	}, &listAggr)

	if err != nil {
		return nil, err
	}

	return listAggr, nil
}
