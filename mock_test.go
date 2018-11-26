package binance_test

import (
	"time"

	"github.com/binance-exchange/go-binance"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ServiceMock) Time() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *ServiceMock) OrderBook(obr binance.OrderBookRequest) (*binance.OrderBook, error) {
	args := m.Called(obr)
	ob, ok := args.Get(0).(*binance.OrderBook)
	if !ok {
		ob = nil
	}
	return ob, args.Error(1)
}

func (m *ServiceMock) AggTrades(atr binance.AggTradesRequest) ([]*binance.AggTrade, error) {
	args := m.Called(atr)
	atc, ok := args.Get(0).([]*binance.AggTrade)
	if !ok {
		atc = nil
	}
	return atc, args.Error(1)
}
func (m *ServiceMock) Klines(kr binance.KlinesRequest) ([]*binance.Kline, error) {
	args := m.Called(kr)
	kc, ok := args.Get(0).([]*binance.Kline)
	if !ok {
		kc = nil
	}
	return kc, args.Error(1)
}
func (m *ServiceMock) Ticker24(tr binance.TickerRequest) (*binance.Ticker24, error) {
	args := m.Called(tr)
	t24, ok := args.Get(0).(*binance.Ticker24)
	if !ok {
		t24 = nil
	}
	return t24, args.Error(1)
}
func (m *ServiceMock) TickerAllPrices() ([]*binance.PriceTicker, error) {
	args := m.Called()
	ptc, ok := args.Get(0).([]*binance.PriceTicker)
	if !ok {
		ptc = nil
	}
	return ptc, args.Error(1)
}
func (m *ServiceMock) TickerAllBooks() ([]*binance.BookTicker, error) {
	args := m.Called()
	btc, ok := args.Get(0).([]*binance.BookTicker)
	if !ok {
		btc = nil
	}
	return btc, args.Error(1)
}
func (m *ServiceMock) NewOrder(or binance.NewOrderRequest) (*binance.ProcessedOrder, error) {
	args := m.Called(or)
	ob, ok := args.Get(0).(*binance.ProcessedOrder)
	if !ok {
		ob = nil
	}
	return ob, args.Error(1)
}
func (m *ServiceMock) NewOrderTest(or binance.NewOrderRequest) error {
	args := m.Called(or)
	return args.Error(0)
}
func (m *ServiceMock) QueryOrder(qor binance.QueryOrderRequest) (*binance.ExecutedOrder, error) {
	args := m.Called(qor)
	eo, ok := args.Get(0).(*binance.ExecutedOrder)
	if !ok {
		eo = nil
	}
	return eo, args.Error(1)
}
func (m *ServiceMock) CancelOrder(cor binance.CancelOrderRequest) (*binance.CanceledOrder, error) {
	args := m.Called(cor)
	co, ok := args.Get(0).(*binance.CanceledOrder)
	if !ok {
		co = nil
	}
	return co, args.Error(1)
}
func (m *ServiceMock) OpenOrders(oor binance.OpenOrdersRequest) ([]*binance.ExecutedOrder, error) {
	args := m.Called(oor)
	eoc, ok := args.Get(0).([]*binance.ExecutedOrder)
	if !ok {
		eoc = nil
	}
	return eoc, args.Error(1)
}
func (m *ServiceMock) AllOrders(aor binance.AllOrdersRequest) ([]*binance.ExecutedOrder, error) {
	args := m.Called(aor)
	eoc, ok := args.Get(0).([]*binance.ExecutedOrder)
	if !ok {
		eoc = nil
	}
	return eoc, args.Error(1)
}
func (m *ServiceMock) Account(ar binance.AccountRequest) (*binance.Account, error) {
	args := m.Called(ar)
	a, ok := args.Get(0).(*binance.Account)
	if !ok {
		a = nil
	}
	return a, args.Error(1)
}
func (m *ServiceMock) MyTrades(mtr binance.MyTradesRequest) ([]*binance.Trade, error) {
	args := m.Called(mtr)
	tc, ok := args.Get(0).([]*binance.Trade)
	if !ok {
		tc = nil
	}
	return tc, args.Error(1)
}
func (m *ServiceMock) Withdraw(wr binance.WithdrawRequest) (*binance.WithdrawResult, error) {
	args := m.Called(wr)
	wres, ok := args.Get(0).(*binance.WithdrawResult)
	if !ok {
		wres = nil
	}
	return wres, args.Error(1)
}
func (m *ServiceMock) DepositHistory(hr binance.HistoryRequest) ([]*binance.Deposit, error) {
	args := m.Called(hr)
	dc, ok := args.Get(0).([]*binance.Deposit)
	if !ok {
		dc = nil
	}
	return dc, args.Error(1)
}
func (m *ServiceMock) WithdrawHistory(hr binance.HistoryRequest) ([]*binance.Withdrawal, error) {
	args := m.Called(hr)
	wc, ok := args.Get(0).([]*binance.Withdrawal)
	if !ok {
		wc = nil
	}
	return wc, args.Error(1)
}
func (m *ServiceMock) StartUserDataStream() (*binance.Stream, error) {
	args := m.Called()
	s, ok := args.Get(0).(*binance.Stream)
	if !ok {
		s = nil
	}
	return s, args.Error(1)
}
func (m *ServiceMock) KeepAliveUserDataStream(s *binance.Stream) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *ServiceMock) CloseUserDataStream(s *binance.Stream) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *ServiceMock) DepthWebsocket(dwr binance.DepthWebsocketRequest) (chan *binance.DepthEvent, chan struct{}, error) {
	args := m.Called(dwr)
	dech, ok := args.Get(0).(chan *binance.DepthEvent)
	if !ok {
		dech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return dech, sch, args.Error(2)
}
func (m *ServiceMock) KlineWebsocket(kwr binance.KlineWebsocketRequest) (chan *binance.KlineEvent, chan struct{}, error) {
	args := m.Called(kwr)
	kech, ok := args.Get(0).(chan *binance.KlineEvent)
	if !ok {
		kech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return kech, sch, args.Error(2)
}
func (m *ServiceMock) TradeWebsocket(twr binance.TradeWebsocketRequest) (chan *binance.AggTradeEvent, chan struct{}, error) {
	args := m.Called(twr)
	atech, ok := args.Get(0).(chan *binance.AggTradeEvent)
	if !ok {
		atech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return atech, sch, args.Error(2)
}
func (m *ServiceMock) UserDataWebsocket(udwr binance.UserDataWebsocketRequest) (chan *binance.AccountEvent, chan struct{}, error) {
	args := m.Called(udwr)
	aech, ok := args.Get(0).(chan *binance.AccountEvent)
	if !ok {
		aech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return aech, sch, args.Error(2)
}
func (m *ServiceMock) ExchangeInfo() (*binance.ExchangeInfo, error) {
	return nil, nil
}
