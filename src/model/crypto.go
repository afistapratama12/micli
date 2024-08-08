package model

// list of market view model
type MarketView struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
	Data      Data   `json:"data"`
}

type Data struct {
	List []DetailData `json:"list"`
}

type DetailData struct {
	Type                int      `json:"type"`
	Symbol              string   `json:"symbol"`
	BaseAsset           string   `json:"baseAsset"`
	BasePrecision       int      `json:"basePrecision"`
	QuoteAsset          string   `json:"quoteAsset"`
	QuotePrecision      int      `json:"quotePrecision"`
	Filters             []Filter `json:"filters"`
	OrderTypes          []string `json:"orderTypes"`
	IcebergEnable       int      `json:"icebergEnable"`
	OcoEnable           int      `json:"ocoEnable"`
	SpotTradingEnable   int      `json:"spotTradingEnable"`
	MarginTradingEnable int      `json:"marginTradingEnable"`
	Permissions         []string `json:"permissions"`
}

type Filter struct {
	FilterType        string  `json:"filterType"`
	MinPrice          string  `json:"minPrice,omitempty"`
	MaxPrice          string  `json:"maxPrice,omitempty"`
	TickSize          string  `json:"tickSize,omitempty"`
	ApplyToMarket     bool    `json:"applyToMarket"`
	MinQty            string  `json:"minQty,omitempty"`
	MaxQty            string  `json:"maxQty,omitempty"`
	StepSize          string  `json:"stepSize,omitempty"`
	Limit             string  `json:"limit,omitempty"`
	BidMultiplierUp   float64 `json:"bidMultiplierUp,omitempty"`
	BidMultiplierDown float64 `json:"bidMultiplierDown,omitempty"`
	AskMultiplierUp   float64 `json:"askMultiplierUp,omitempty"`
	AskMultiplierDown float64 `json:"askMultiplierDown,omitempty"`
	AvgPriceMins      string  `json:"avgPriceMins,omitempty"`
	MinNotional       string  `json:"minNotional,omitempty"`
	MaxNotional       string  `json:"maxNotional,omitempty"`
	MaxNumAlgoOrders  string  `json:"maxNumAlgoOrders,omitempty"`
}

type ViewPairMarket struct {
	USD    string
	IDR    string
	BTC    string
	BNB    string
	Others []string
}

type AggrTradeData struct {
	AggrId           int64  `json:"a"`
	Price            string `json:"p"`
	Qty              string `json:"q"`
	FirstId          int64  `json:"f"`
	LastId           int64  `json:"l"`
	Timestamp        int64  `json:"T"` // in milisecond
	IsMaker          bool   `json:"m"`
	IsBestPriceMatch bool   `json:"M"`
}

// {
// 	"a": 2468121,
// 	"p": "16100.00",
// 	"q": "38.50000000",
// 	"f": 4663729,
// 	"l": 4663729,
// 	"T": 1723082902996,
// 	"m": false
// }

type DepthData struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type BaseWSData struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
}

type WSDepthData struct {
	BaseWSData
	UpdateId int64      `json:"u"`
	FirstId  int64      `json:"U"`
	Bids     [][]string `json:"b"`
	Asks     [][]string `json:"a"`
}

func (w WSDepthData) ToDepthData() DepthData {
	return DepthData{
		LastUpdateId: w.UpdateId,
		Bids:         w.Bids,
		Asks:         w.Asks,
	}
}

type WSAggrTradeData struct {
	BaseWSData
	AggrTradeId      int64  `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeId     int64  `json:"f"`
	LastTradeId      int64  `json:"l"`
	TradeTime        int64  `json:"T"`
	IsBuyerMaker     bool   `json:"m"`
	IsBestPriceMatch bool   `json:"M"`
}

func (w WSAggrTradeData) ToAggrTradeData() AggrTradeData {
	return AggrTradeData{
		AggrId:           w.AggrTradeId,
		Price:            w.Price,
		Qty:              w.Quantity,
		FirstId:          w.FirstTradeId,
		LastId:           w.LastTradeId,
		Timestamp:        w.TradeTime,
		IsMaker:          w.IsBuyerMaker,
		IsBestPriceMatch: w.IsBestPriceMatch,
	}
}

type CacheData struct {
	ListPair []string
}
