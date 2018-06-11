package binance_test

import (
	"testing"
	"time"

	"github.com/bwmuller/go-binance"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	nor := binance.NewOrderRequest{
		Symbol:      "BNBETH",
		Quantity:    1,
		Price:       999,
		Side:        binance.SideSell,
		TimeInForce: binance.GTC,
		Type:        binance.TypeLimit,
		Timestamp:   time.Now(),
	}

	po := &binance.ProcessedOrder{
		Symbol:        nor.Symbol,
		OrderID:       int64(123),
		ClientOrderID: "clientOrderID",
		TransactTime:  time.Now(),
	}
	binanceService.On("NewOrder", nor).Return(po, nil)
	po_r, err := b.NewOrder(nor)
	assert.Nil(t, err)
	assert.Equal(t, po, po_r)
}

func TestNewOrderTest(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	nor := binance.NewOrderRequest{
		Symbol:      "BNBETH",
		Quantity:    1,
		Price:       999,
		Side:        binance.SideSell,
		TimeInForce: binance.GTC,
		Type:        binance.TypeLimit,
		Timestamp:   time.Now(),
	}

	binanceService.On("NewOrderTest", nor).Return(nil)
	err := b.NewOrderTest(nor)
	assert.Nil(t, err)
}

func TestQueryOrder(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	qor := binance.QueryOrderRequest{
		Symbol:            "BNBETH",
		OrderID:           int64(123),
		OrigClientOrderID: "clientOrderID",
		RecvWindow:        1 * time.Second,
		Timestamp:         time.Now(),
	}
	eo := &binance.ExecutedOrder{
		Symbol:        "BNBETH",
		OrderID:       123,
		ClientOrderID: "clientOrderID",
		Price:         10.23,
		OrigQty:       10.00,
		ExecutedQty:   4.44,
		Status:        binance.StatusPartiallyFilled,
		TimeInForce:   binance.GTC,
		Type:          binance.TypeLimit,
		Side:          binance.SideBuy,
		Time:          time.Now(),
	}

	binanceService.On("QueryOrder", qor).Return(eo, nil)
	eo_r, err := b.QueryOrder(qor)
	assert.Nil(t, err)
	assert.Equal(t, eo, eo_r)
}

func TestCancelOrder(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	cor := binance.CancelOrderRequest{
		Symbol:            "BNBETH",
		OrderID:           int64(123),
		OrigClientOrderID: "clientOrderID",
		RecvWindow:        1 * time.Second,
		Timestamp:         time.Now(),
	}
	co := &binance.CanceledOrder{
		Symbol:        "BNBETH",
		OrderID:       123,
		ClientOrderID: "clientOrderID",
	}

	binanceService.On("CancelOrder", cor).Return(co, nil)
	co_r, err := b.CancelOrder(cor)
	assert.Nil(t, err)
	assert.Equal(t, co, co_r)
}

func TestOpenOrders(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	oor := binance.OpenOrdersRequest{
		Symbol:     "BNBETH",
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	ooc := []*binance.ExecutedOrder{}

	binanceService.On("OpenOrders", oor).Return(ooc, nil)
	ooc_r, err := b.OpenOrders(oor)
	assert.Nil(t, err)
	assert.Equal(t, ooc, ooc_r)
}

func TestAllOrders(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	aor := binance.AllOrdersRequest{
		Symbol:     "BNBETH",
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
		Limit:      10,
	}
	aoc := []*binance.ExecutedOrder{}

	binanceService.On("AllOrders", aor).Return(aoc, nil)
	aoc_r, err := b.AllOrders(aor)
	assert.Nil(t, err)
	assert.Equal(t, aoc, aoc_r)
}

func TestAccount(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	ar := binance.AccountRequest{
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	a := &binance.Account{}

	binanceService.On("Account", ar).Return(a, nil)
	a_r, err := b.Account(ar)
	assert.Nil(t, err)
	assert.Equal(t, a, a_r)
}

func TestMyTrades(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	mtr := binance.MyTradesRequest{
		Symbol:     "BNBETH",
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	tc := []*binance.Trade{}

	binanceService.On("MyTrades", mtr).Return(tc, nil)
	tc_r, err := b.MyTrades(mtr)
	assert.Nil(t, err)
	assert.Equal(t, tc, tc_r)
}

func TestWithdraw(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	wr := binance.WithdrawRequest{
		Asset:      "ETH",
		Address:    "0x1234",
		Amount:     1.23,
		Name:       "My wallet",
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	wres := &binance.WithdrawResult{
		Success: true,
	}

	binanceService.On("Withdraw", wr).Return(wres, nil)
	wres_r, err := b.Withdraw(wr)
	assert.Nil(t, err)
	assert.Equal(t, wres, wres_r)
}

func TestDepositHistory(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	hr := binance.HistoryRequest{
		Asset:      "ETH",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	dhc := []*binance.Deposit{}

	binanceService.On("DepositHistory", hr).Return(dhc, nil)
	dhc_r, err := b.DepositHistory(hr)
	assert.Nil(t, err)
	assert.Equal(t, dhc, dhc_r)
}

func TestWithdrawHistory(t *testing.T) {
	binanceService := &ServiceMock{}
	b := binance.NewBinance(binanceService)

	hr := binance.HistoryRequest{
		Asset:      "ETH",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
		RecvWindow: 1 * time.Second,
		Timestamp:  time.Now(),
	}
	dhc := []*binance.Withdrawal{}

	binanceService.On("WithdrawHistory", hr).Return(dhc, nil)
	dhc_r, err := b.WithdrawHistory(hr)
	assert.Nil(t, err)
	assert.Equal(t, dhc, dhc_r)
}
