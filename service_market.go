package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func (as *apiService) Ping() error {
	params := make(map[string]string)
	response, err := as.request("GET", "api/v1/ping", params, false, false)
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", response.StatusCode)
	return nil
}

func (as *apiService) Time() (time.Time, error) {
	params := make(map[string]string)
	res, err := as.request("GET", "api/v1/time", params, false, false)
	if err != nil {
		return time.Time{}, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()
	var rawTime struct {
		ServerTime string `json:"serverTime"`
	}
	if err := json.Unmarshal(textRes, &rawTime); err != nil {
		return time.Time{}, errors.Wrap(err, "timeResponse unmarshal failed")
	}
	t, err := timeFromUnixTimestampFloat(rawTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (as *apiService) OrderBook(obr OrderBookRequest) (*OrderBook, error) {
	params := make(map[string]string)
	params["symbol"] = obr.Symbol
	if obr.Limit != 0 {
		params["limit"] = strconv.Itoa(obr.Limit)
	}
	res, err := as.request("GET", "api/v1/depth", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Time")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawBook := &struct {
		LastUpdateID int             `json:"lastUpdateId"`
		Bids         [][]interface{} `json:"bids"`
		Asks         [][]interface{} `json:"asks"`
	}{}
	if err := json.Unmarshal(textRes, rawBook); err != nil {
		return nil, errors.Wrap(err, "timeResponse unmarshal failed")
	}

	ob := &OrderBook{
		LastUpdateID: rawBook.LastUpdateID,
	}
	extractOrder := func(rawPrice, rawQuantity interface{}) (*Order, error) {
		price, err := floatFromString(rawPrice)
		if err != nil {
			return nil, err
		}
		quantity, err := floatFromString(rawQuantity)
		if err != nil {
			return nil, err
		}
		return &Order{
			Price:    price,
			Quantity: quantity,
		}, nil
	}
	for _, bid := range rawBook.Bids {
		order, err := extractOrder(bid[0], bid[1])
		if err != nil {
			return nil, err
		}
		ob.Bids = append(ob.Bids, order)
	}
	for _, ask := range rawBook.Asks {
		order, err := extractOrder(ask[0], ask[1])
		if err != nil {
			return nil, err
		}
		ob.Asks = append(ob.Asks, order)
	}

	return ob, nil
}

func (as *apiService) AggTrades(atr AggTradesRequest) ([]*AggTrade, error) {
	params := make(map[string]string)
	params["symbol"] = atr.Symbol
	if atr.FromID != 0 {
		params["fromId"] = strconv.FormatInt(atr.FromID, 10)
	}
	if atr.StartTime != 0 {
		params["startTime"] = strconv.FormatInt(atr.StartTime, 10)
	}
	if atr.EndTime != 0 {
		params["endTime"] = strconv.FormatInt(atr.EndTime, 10)
	}
	if atr.Limit != 0 {
		params["limit"] = strconv.Itoa(atr.Limit)
	}

	res, err := as.request("GET", "api/v1/aggTrades", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from AggTrades")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawAggTrades := []struct {
		ID             int    `json:"a"`
		Price          string `json:"p"`
		Quantity       string `json:"q"`
		FirstTradeID   int    `json:"f"`
		LastTradeID    int    `json:"l"`
		Timestamp      int64  `json:"T"`
		BuyerMaker     bool   `json:"m"`
		BestPriceMatch bool   `json:"M"`
	}{}
	if err := json.Unmarshal(textRes, &rawAggTrades); err != nil {
		return nil, errors.Wrap(err, "aggTrades unmarshal failed")
	}
	aggTrades := []*AggTrade{}
	for _, rawTrade := range rawAggTrades {
		price, err := floatFromString(rawTrade.Price)
		if err != nil {
			return nil, err
		}
		quantity, err := floatFromString(rawTrade.Quantity)
		if err != nil {
			return nil, err
		}
		t := time.Unix(0, rawTrade.Timestamp*int64(time.Millisecond))

		aggTrades = append(aggTrades, &AggTrade{
			ID:             rawTrade.ID,
			Price:          price,
			Quantity:       quantity,
			FirstTradeID:   rawTrade.FirstTradeID,
			LastTradeID:    rawTrade.LastTradeID,
			Timestamp:      t,
			BuyerMaker:     rawTrade.BuyerMaker,
			BestPriceMatch: rawTrade.BestPriceMatch,
		})
	}
	return aggTrades, nil
}

func (as *apiService) RecentTrades(rtr RecentTradesRequest) ([]*RecentTrades, error) {
	params := make(map[string]string)
	params["symbol"] = rtr.Symbol
	if rtr.Limit != 0 {
		params["limit"] = strconv.Itoa(rtr.Limit)
	}

	res, err := as.request("GET", "api/v1/trades", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from RecentTrades")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	var rawRecentTrades []struct {
		ID              int64   `json:"id"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		Time            float64 `json:"time"`
		IsBuyerMaker    bool    `json:"isBuyerMaker"`
		IsBestMatch     bool    `json:"isBestMatch"`
	}

	if err := json.Unmarshal(textRes, &rawRecentTrades); err != nil {
		return nil, errors.Wrap(err, "rawRecentTrades unmarshal failed")
	}

	var recentTrades []*RecentTrades
	for _, rrt := range rawRecentTrades {
		price, err := floatFromString(rrt.Price)
		if err != nil {
			return nil, err
		}
		qty, err := floatFromString(rrt.Qty)
		if err != nil {
			return nil, err
		}

		t, err := timeFromUnixTimestampFloat(rrt.Time)
		if err != nil {
			return nil, err
		}
		recentTrades = append(recentTrades, &RecentTrades{
			ID:              rrt.ID,
			Price:           price,
			Qty:             qty,
			Time:            t,
			IsBuyerMaker:    rrt.IsBuyerMaker,
			IsBestMatch:     rrt.IsBestMatch,
		})
	}
	return recentTrades, nil
}

func (as *apiService) Klines(kr KlinesRequest) ([]*Kline, error) {
	params := make(map[string]string)
	params["symbol"] = kr.Symbol
	params["interval"] = string(kr.Interval)
	if kr.Limit != 0 {
		params["limit"] = strconv.Itoa(kr.Limit)
	}
	if kr.StartTime != 0 {
		params["startTime"] = strconv.FormatInt(kr.StartTime, 10)
	}
	if kr.EndTime != 0 {
		params["endTime"] = strconv.FormatInt(kr.EndTime, 10)
	}

	res, err := as.request("GET", "api/v1/klines", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Klines")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawKlines := [][]interface{}{}
	if err := json.Unmarshal(textRes, &rawKlines); err != nil {
		return nil, errors.Wrap(err, "rawKlines unmarshal failed")
	}
	klines := []*Kline{}
	for _, k := range rawKlines {
		ot, err := timeFromUnixTimestampFloat(k[0])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.OpenTime")
		}
		open, err := floatFromString(k[1])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Open")
		}
		high, err := floatFromString(k[2])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.High")
		}
		low, err := floatFromString(k[3])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Low")
		}
		cls, err := floatFromString(k[4])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Close")
		}
		volume, err := floatFromString(k[5])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.Volume")
		}
		ct, err := timeFromUnixTimestampFloat(k[6])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.CloseTime")
		}
		qav, err := floatFromString(k[7])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.QuoteAssetVolume")
		}
		not, ok := k[8].(float64)
		if !ok {
			return nil, errors.Wrap(err, "cannot parse Kline.NumberOfTrades")
		}
		tbbav, err := floatFromString(k[9])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyBaseAssetVolume")
		}
		tbqav, err := floatFromString(k[10])
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse Kline.TakerBuyQuoteAssetVolume")
		}
		klines = append(klines, &Kline{
			OpenTime:                 ot,
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    cls,
			Volume:                   volume,
			CloseTime:                ct,
			QuoteAssetVolume:         qav,
			NumberOfTrades:           int(not),
			TakerBuyBaseAssetVolume:  tbbav,
			TakerBuyQuoteAssetVolume: tbqav,
		})
	}
	return klines, nil
}

func (as *apiService) Ticker24(tr TickerRequest) (*Ticker24, error) {
	params := make(map[string]string)
	params["symbol"] = tr.Symbol

	res, err := as.request("GET", "api/v1/ticker/24hr", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawTicker24 := struct {
		PriceChange        string  `json:"priceChange"`
		PriceChangePercent string  `json:"priceChangePercent"`
		WeightedAvgPrice   string  `json:"weightedAvgPrice"`
		PrevClosePrice     string  `json:"prevClosePrice"`
		LastPrice          string  `json:"lastPrice"`
		BidPrice           string  `json:"bidPrice"`
		AskPrice           string  `json:"askPrice"`
		OpenPrice          string  `json:"openPrice"`
		HighPrice          string  `json:"highPrice"`
		LowPrice           string  `json:"lowPrice"`
		Volume             string  `json:"volume"`
		OpenTime           float64 `json:"openTime"`
		CloseTime          float64 `json:"closeTime"`
		FirstID            int
		LastID             int
		Count              int
	}{}
	if err := json.Unmarshal(textRes, &rawTicker24); err != nil {
		return nil, errors.Wrap(err, "rawTicker24 unmarshal failed")
	}

	pc, err := strconv.ParseFloat(rawTicker24.PriceChange, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PriceChange")
	}
	pcPercent, err := strconv.ParseFloat(rawTicker24.PriceChangePercent, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PriceChangePercent")
	}
	wap, err := strconv.ParseFloat(rawTicker24.WeightedAvgPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.WeightedAvgPrice")
	}
	pcp, err := strconv.ParseFloat(rawTicker24.PrevClosePrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.PrevClosePrice")
	}
	lastPrice, err := strconv.ParseFloat(rawTicker24.LastPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.LastPrice")
	}
	bp, err := strconv.ParseFloat(rawTicker24.BidPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.BidPrice")
	}
	ap, err := strconv.ParseFloat(rawTicker24.AskPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.AskPrice")
	}
	op, err := strconv.ParseFloat(rawTicker24.OpenPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.OpenPrice")
	}
	hp, err := strconv.ParseFloat(rawTicker24.HighPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.HighPrice")
	}
	lowPrice, err := strconv.ParseFloat(rawTicker24.LowPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.LowPrice")
	}
	vol, err := strconv.ParseFloat(rawTicker24.Volume, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.Volume")
	}
	ot, err := timeFromUnixTimestampFloat(rawTicker24.OpenTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.OpenTime")
	}
	ct, err := timeFromUnixTimestampFloat(rawTicker24.CloseTime)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Ticker24.CloseTime")
	}
	t24 := &Ticker24{
		PriceChange:        pc,
		PriceChangePercent: pcPercent,
		WeightedAvgPrice:   wap,
		PrevClosePrice:     pcp,
		LastPrice:          lastPrice,
		BidPrice:           bp,
		AskPrice:           ap,
		OpenPrice:          op,
		HighPrice:          hp,
		LowPrice:           lowPrice,
		Volume:             vol,
		OpenTime:           ot,
		CloseTime:          ct,
		FirstID:            rawTicker24.FirstID,
		LastID:             rawTicker24.LastID,
		Count:              rawTicker24.Count,
	}
	return t24, nil
}

func (as *apiService) TickerAllPrices() ([]*PriceTicker, error) {
	params := make(map[string]string)

	res, err := as.request("GET", "api/v1/ticker/allPrices", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		as.handleError(textRes)
	}

	rawTickerAllPrices := []struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}{}
	if err := json.Unmarshal(textRes, &rawTickerAllPrices); err != nil {
		return nil, errors.Wrap(err, "rawTickerAllPrices unmarshal failed")
	}

	var tpc []*PriceTicker
	for _, rawTickerPrice := range rawTickerAllPrices {
		p, err := strconv.ParseFloat(rawTickerPrice.Price, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerAllPrices.Price")
		}
		tpc = append(tpc, &PriceTicker{
			Symbol: rawTickerPrice.Symbol,
			Price:  p,
		})
	}
	return tpc, nil
}

func (as *apiService) TickerAllBooks() ([]*BookTicker, error) {
	params := make(map[string]string)

	res, err := as.request("GET", "api/v1/ticker/allBookTickers", params, false, false)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/allBookTickers")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawBookTickers := []struct {
		Symbol   string `json:"symbol"`
		BidPrice string `json:"bidPrice"`
		BidQty   string `json:"bidQty"`
		AskPrice string `json:"askPrice"`
		AskQty   string `json:"askQty"`
	}{}
	if err := json.Unmarshal(textRes, &rawBookTickers); err != nil {
		return nil, errors.Wrap(err, "rawBookTickers unmarshal failed")
	}

	var btc []*BookTicker
	for _, rawBookTicker := range rawBookTickers {
		bp, err := strconv.ParseFloat(rawBookTicker.BidPrice, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.BidPrice")
		}
		bqty, err := strconv.ParseFloat(rawBookTicker.BidQty, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.BidQty")
		}
		ap, err := strconv.ParseFloat(rawBookTicker.AskPrice, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.AskPrice")
		}
		aqty, err := strconv.ParseFloat(rawBookTicker.AskQty, 64)
		if err != nil {
			return nil, errors.Wrap(err, "cannot parse TickerBookTickers.AskQty")
		}
		btc = append(btc, &BookTicker{
			Symbol:   rawBookTicker.Symbol,
			BidPrice: bp,
			BidQty:   bqty,
			AskPrice: ap,
			AskQty:   aqty,
		})
	}
	return btc, nil
}
