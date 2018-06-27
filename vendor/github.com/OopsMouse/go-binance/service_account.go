package binance

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

type rawExecutedOrder struct {
	Symbol        string  `json:"symbol"`
	OrderID       int64   `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	Price         string  `json:"price"`
	OrigQty       string  `json:"origQty"`
	ExecutedQty   string  `json:"executedQty"`
	Status        string  `json:"status"`
	TimeInForce   string  `json:"timeInForce"`
	Type          string  `json:"type"`
	Side          string  `json:"side"`
	StopPrice     string  `json:"stopPrice"`
	IcebergQty    string  `json:"icebergQty"`
	Time          float64 `json:"time"`
}

func (as *apiService) NewOrder(or NewOrderRequest) (*ProcessedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)
	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 8, 64)
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', 8, 64)
	}
	if string(or.TimeInForce) != "" {
		params["timeInForce"] = string(or.TimeInForce)
	}
	if or.Price > 0 {
		params["price"] = strconv.FormatFloat(or.Price, 'f', 8, 64)
	}

	res, err := as.request("POST", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := struct {
		Symbol        string  `json:"symbol"`
		OrderID       int64   `json:"orderId"`
		ClientOrderID string  `json:"clientOrderId"`
		TransactTime  float64 `json:"transactTime"`
	}{}
	if err := json.Unmarshal(textRes, &rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	t, err := timeFromUnixTimestampFloat(rawOrder.TransactTime)
	if err != nil {
		return nil, err
	}

	return &ProcessedOrder{
		Symbol:        rawOrder.Symbol,
		OrderID:       rawOrder.OrderID,
		ClientOrderID: rawOrder.ClientOrderID,
		TransactTime:  t,
	}, nil
}

func (as *apiService) NewOrderTest(or NewOrderRequest) error {
	params := make(map[string]string)
	params["symbol"] = or.Symbol
	params["side"] = string(or.Side)
	params["type"] = string(or.Type)
	params["quantity"] = strconv.FormatFloat(or.Quantity, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(or.Timestamp), 10)
	if or.NewClientOrderID != "" {
		params["newClientOrderId"] = or.NewClientOrderID
	}
	if or.StopPrice != 0 {
		params["stopPrice"] = strconv.FormatFloat(or.StopPrice, 'f', 8, 64)
	}
	if or.IcebergQty != 0 {
		params["icebergQty"] = strconv.FormatFloat(or.IcebergQty, 'f', 8, 64)
	}
	if string(or.TimeInForce) != "" {
		params["timeInForce"] = string(or.TimeInForce)
	}
	if or.Price > 0 {
		params["price"] = strconv.FormatFloat(or.Price, 'f', 8, 64)
	}

	res, err := as.request("POST", "api/v3/order/test", params, true, true)
	if err != nil {
		return err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read response from Ticker/24hr")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return as.handleError(textRes)
	}
	return nil
}

func (as *apiService) QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = qor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(qor.Timestamp), 10)
	if qor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(qor.OrderID, 10)
	}
	if qor.OrigClientOrderID != "" {
		params["origClientOrderId"] = qor.OrigClientOrderID
	}
	if qor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(qor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrder := &rawExecutedOrder{}
	if err := json.Unmarshal(textRes, rawOrder); err != nil {
		return nil, errors.Wrap(err, "rawOrder unmarshal failed")
	}

	eo, err := executedOrderFromRaw(rawOrder)
	if err != nil {
		return nil, err
	}
	return eo, nil
}

func (as *apiService) CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	params := make(map[string]string)
	params["symbol"] = cor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(cor.Timestamp), 10)
	if cor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(cor.OrderID, 10)
	}
	if cor.OrigClientOrderID != "" {
		params["origClientOrderId"] = cor.OrigClientOrderID
	}
	if cor.NewClientOrderID != "" {
		params["newClientOrderId"] = cor.NewClientOrderID
	}
	if cor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(cor.RecvWindow), 10)
	}

	res, err := as.request("DELETE", "api/v3/order", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from order.delete")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawCanceledOrder := struct {
		Symbol            string `json:"symbol"`
		OrigClientOrderID string `json:"origClientOrderId"`
		OrderID           int64  `json:"orderId"`
		ClientOrderID     string `json:"clientOrderId"`
	}{}
	if err := json.Unmarshal(textRes, &rawCanceledOrder); err != nil {
		return nil, errors.Wrap(err, "cancelOrder unmarshal failed")
	}

	return &CanceledOrder{
		Symbol:            rawCanceledOrder.Symbol,
		OrigClientOrderID: rawCanceledOrder.OrigClientOrderID,
		OrderID:           rawCanceledOrder.OrderID,
		ClientOrderID:     rawCanceledOrder.ClientOrderID,
	}, nil
}

func (as *apiService) OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = oor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(oor.Timestamp), 10)
	if oor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(oor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/openOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from openOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "openOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	params := make(map[string]string)
	params["symbol"] = aor.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(aor.Timestamp), 10)
	if aor.OrderID != 0 {
		params["orderId"] = strconv.FormatInt(aor.OrderID, 10)
	}
	if aor.Limit != 0 {
		params["limit"] = strconv.Itoa(aor.Limit)
	}
	if aor.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(aor.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/allOrders", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from allOrders.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawOrders := []*rawExecutedOrder{}
	if err := json.Unmarshal(textRes, &rawOrders); err != nil {
		return nil, errors.Wrap(err, "allOrders unmarshal failed")
	}

	var eoc []*ExecutedOrder
	for _, rawOrder := range rawOrders {
		eo, err := executedOrderFromRaw(rawOrder)
		if err != nil {
			return nil, err
		}
		eoc = append(eoc, eo)
	}

	return eoc, nil
}

func (as *apiService) Account(ar AccountRequest) (*Account, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(ar.Timestamp), 10)
	if ar.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(ar.RecvWindow), 10)
	}

	res, err := as.request("GET", "api/v3/account", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from account.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawAccount := struct {
		MakerCommision   int64 `json:"makerCommision"`
		TakerCommission  int64 `json:"takerCommission"`
		BuyerCommission  int64 `json:"buyerCommission"`
		SellerCommission int64 `json:"sellerCommission"`
		CanTrade         bool  `json:"canTrade"`
		CanWithdraw      bool  `json:"canWithdraw"`
		CanDeposit       bool  `json:"canDeposit"`
		Balances         []struct {
			Asset  string `json:"asset"`
			Free   string `json:"free"`
			Locked string `json:"locked"`
		}
	}{}
	if err := json.Unmarshal(textRes, &rawAccount); err != nil {
		return nil, errors.Wrap(err, "rawAccount unmarshal failed")
	}

	acc := &Account{
		MakerCommision:  rawAccount.MakerCommision,
		TakerCommision:  rawAccount.TakerCommission,
		BuyerCommision:  rawAccount.BuyerCommission,
		SellerCommision: rawAccount.SellerCommission,
		CanTrade:        rawAccount.CanTrade,
		CanWithdraw:     rawAccount.CanWithdraw,
		CanDeposit:      rawAccount.CanDeposit,
	}
	for _, b := range rawAccount.Balances {
		f, err := floatFromString(b.Free)
		if err != nil {
			return nil, err
		}
		l, err := floatFromString(b.Locked)
		if err != nil {
			return nil, err
		}
		acc.Balances = append(acc.Balances, &Balance{
			Asset:  b.Asset,
			Free:   f,
			Locked: l,
		})
	}

	return acc, nil
}

func (as *apiService) MyTrades(mtr MyTradesRequest) ([]*Trade, error) {
	params := make(map[string]string)
	params["symbol"] = mtr.Symbol
	params["timestamp"] = strconv.FormatInt(unixMillis(mtr.Timestamp), 10)
	if mtr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(mtr.RecvWindow), 10)
	}
	if mtr.FromID != 0 {
		params["orderId"] = strconv.FormatInt(mtr.FromID, 10)
	}
	if mtr.Limit != 0 {
		params["limit"] = strconv.Itoa(mtr.Limit)
	}

	res, err := as.request("GET", "api/v3/myTrades", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from myTrades.get")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawTrades := []struct {
		ID              int64   `json:"id"`
		Price           string  `json:"price"`
		Qty             string  `json:"qty"`
		Commission      string  `json:"commission"`
		CommissionAsset string  `json:"commissionAsset"`
		Time            float64 `json:"time"`
		IsBuyer         bool    `json:"isBuyer"`
		IsMaker         bool    `json:"isMaker"`
		IsBestMatch     bool    `json:"isBestMatch"`
	}{}
	if err := json.Unmarshal(textRes, &rawTrades); err != nil {
		return nil, errors.Wrap(err, "rawTrades unmarshal failed")
	}

	var tc []*Trade
	for _, rt := range rawTrades {
		price, err := floatFromString(rt.Price)
		if err != nil {
			return nil, err
		}
		qty, err := floatFromString(rt.Qty)
		if err != nil {
			return nil, err
		}
		commission, err := floatFromString(rt.Commission)
		if err != nil {
			return nil, err
		}
		t, err := timeFromUnixTimestampFloat(rt.Time)
		if err != nil {
			return nil, err
		}
		tc = append(tc, &Trade{
			ID:              rt.ID,
			Price:           price,
			Qty:             qty,
			Commission:      commission,
			CommissionAsset: rt.CommissionAsset,
			Time:            t,
			IsBuyer:         rt.IsBuyer,
			IsMaker:         rt.IsMaker,
			IsBestMatch:     rt.IsBestMatch,
		})
	}
	return tc, nil
}

func (as *apiService) Withdraw(wr WithdrawRequest) (*WithdrawResult, error) {
	params := make(map[string]string)
	params["asset"] = wr.Asset
	params["address"] = wr.Address
	params["amount"] = strconv.FormatFloat(wr.Amount, 'f', 8, 64)
	params["timestamp"] = strconv.FormatInt(unixMillis(wr.Timestamp), 10)
	if wr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(wr.RecvWindow), 10)
	}
	if wr.Name != "" {
		params["name"] = wr.Name
	}

	res, err := as.request("POST", "wapi/v1/withdraw.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from withdraw.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawResult := struct {
		Msg     string `json:"msg"`
		Success bool   `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawResult); err != nil {
		return nil, errors.Wrap(err, "rawTrades unmarshal failed")
	}

	return &WithdrawResult{
		Msg:     rawResult.Msg,
		Success: rawResult.Success,
	}, nil
}
func (as *apiService) DepositHistory(hr HistoryRequest) ([]*Deposit, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(hr.Timestamp), 10)
	if hr.Asset != "" {
		params["asset"] = hr.Asset
	}
	if hr.Status != nil {
		params["status"] = strconv.Itoa(*hr.Status)
	}
	if !hr.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.StartTime), 10)
	}
	if !hr.EndTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.EndTime), 10)
	}
	if hr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(hr.RecvWindow), 10)
	}

	res, err := as.request("POST", "wapi/v1/getDepositHistory.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from depositHistory.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawDepositHistory := struct {
		DepositList []struct {
			InsertTime float64 `json:"insertTime"`
			Amount     float64 `json:"amount"`
			Asset      string  `json:"asset"`
			Status     int     `json:"status"`
		}
		Success bool `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawDepositHistory); err != nil {
		return nil, errors.Wrap(err, "rawDepositHistory unmarshal failed")
	}

	var dc []*Deposit
	for _, d := range rawDepositHistory.DepositList {
		t, err := timeFromUnixTimestampFloat(d.InsertTime)
		if err != nil {
			return nil, err
		}
		dc = append(dc, &Deposit{
			InsertTime: t,
			Amount:     d.Amount,
			Asset:      d.Asset,
			Status:     d.Status,
		})
	}

	return dc, nil
}
func (as *apiService) WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error) {
	params := make(map[string]string)
	params["timestamp"] = strconv.FormatInt(unixMillis(hr.Timestamp), 10)
	if hr.Asset != "" {
		params["asset"] = hr.Asset
	}
	if hr.Status != nil {
		params["status"] = strconv.Itoa(*hr.Status)
	}
	if !hr.StartTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.StartTime), 10)
	}
	if !hr.EndTime.IsZero() {
		params["startTime"] = strconv.FormatInt(unixMillis(hr.EndTime), 10)
	}
	if hr.RecvWindow != 0 {
		params["recvWindow"] = strconv.FormatInt(recvWindow(hr.RecvWindow), 10)
	}

	res, err := as.request("POST", "wapi/v1/getWithdrawHistory.html", params, true, true)
	if err != nil {
		return nil, err
	}
	textRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response from withdrawHistory.post")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, as.handleError(textRes)
	}

	rawWithdrawHistory := struct {
		WithdrawList []struct {
			Amount    float64 `json:"amount"`
			Address   string  `json:"address"`
			TxID      string  `json:"txId"`
			Asset     string  `json:"asset"`
			ApplyTime float64 `json:"insertTime"`
			Status    int     `json:"status"`
		}
		Success bool `json:"success"`
	}{}
	if err := json.Unmarshal(textRes, &rawWithdrawHistory); err != nil {
		return nil, errors.Wrap(err, "rawWithdrawHistory unmarshal failed")
	}

	var wc []*Withdrawal
	for _, w := range rawWithdrawHistory.WithdrawList {
		t, err := timeFromUnixTimestampFloat(w.ApplyTime)
		if err != nil {
			return nil, err
		}
		wc = append(wc, &Withdrawal{
			Amount:    w.Amount,
			Address:   w.Address,
			TxID:      w.TxID,
			Asset:     w.Asset,
			ApplyTime: t,
			Status:    w.Status,
		})
	}

	return wc, nil
}

func executedOrderFromRaw(reo *rawExecutedOrder) (*ExecutedOrder, error) {
	price, err := strconv.ParseFloat(reo.Price, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.CloseTime")
	}
	origQty, err := strconv.ParseFloat(reo.OrigQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.OrigQty")
	}
	execQty, err := strconv.ParseFloat(reo.ExecutedQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.ExecutedQty")
	}
	stopPrice, err := strconv.ParseFloat(reo.StopPrice, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.StopPrice")
	}
	icebergQty, err := strconv.ParseFloat(reo.IcebergQty, 64)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.IcebergQty")
	}
	t, err := timeFromUnixTimestampFloat(reo.Time)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse Order.CloseTime")
	}

	return &ExecutedOrder{
		Symbol:        reo.Symbol,
		OrderID:       reo.OrderID,
		ClientOrderID: reo.ClientOrderID,
		Price:         price,
		OrigQty:       origQty,
		ExecutedQty:   execQty,
		Status:        OrderStatus(reo.Status),
		TimeInForce:   TimeInForce(reo.TimeInForce),
		Type:          OrderType(reo.Type),
		Side:          OrderSide(reo.Side),
		StopPrice:     stopPrice,
		IcebergQty:    icebergQty,
		Time:          t,
	}, nil
}
