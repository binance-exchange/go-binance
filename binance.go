package binance

import (
	"fmt"
	"time"
)

// Binance is wrapper for Binance API.
//
// Read web documentation for more endpoints descriptions and list of
// mandatory and optional params. Wrapper is not responsible for client-side
// validation and only sends requests further.
//
// For each API-defined enum there's a special type and list of defined
// enum values to be used.
type Binance interface {
	// Ping tests connectivity.
	Ping() error
	// Time returns server time.
	Time() (time.Time, error)
	// OrderBook returns list of orders.
	OrderBook(obr OrderBookRequest) (*OrderBook, error)
	// AggTrades returns compressed/aggregate list of trades.
	AggTrades(atr AggTradesRequest) ([]*AggTrade, error)
	// Klines returns klines/candlestick data.
	Klines(kr KlinesRequest) ([]*Kline, error)
	// Ticker24 returns 24hr price change statistics.
	Ticker24(tr TickerRequest) (*Ticker24, error)
	// TickerAllPrices returns ticker data for symbols.
	TickerAllPrices() ([]*PriceTicker, error)
	// TickerAllBooks returns tickers for all books.
	TickerAllBooks() ([]*BookTicker, error)

	// NewOrder places new order and returns ProcessedOrder.
	NewOrder(nor NewOrderRequest) (*ProcessedOrder, error)
	// NewOrder places testing order.
	NewOrderTest(nor NewOrderRequest) error
	// QueryOrder returns data about existing order.
	QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error)
	// CancelOrder cancels order.
	CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error)
	// OpenOrders returns list of open orders.
	OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error)
	// AllOrders returns list of all previous orders.
	AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error)
	// RecentTrades returns all recent trades
	RecentTrades(hr RecentTradesRequest) ([]*RecentTrades, error)

	// Account returns account data.
	Account(ar AccountRequest) (*Account, error)
	// MyTrades list user's trades.
	MyTrades(mtr MyTradesRequest) ([]*Trade, error)
	// Withdraw executes withdrawal.
	Withdraw(wr WithdrawRequest) (*WithdrawResult, error)
	// DepositHistory lists deposit data.
	DepositHistory(hr HistoryRequest) ([]*Deposit, error)
	// WithdrawHistory lists withdraw data.
	WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error)

	// StartUserDataStream starts stream and returns Stream with ListenKey.
	StartUserDataStream() (*Stream, error)
	// KeepAliveUserDataStream prolongs stream livespan.
	KeepAliveUserDataStream(s *Stream) error
	// CloseUserDataStream closes opened stream.
	CloseUserDataStream(s *Stream) error

	DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error)
}

type binance struct {
	Service Service
}

// Error represents Binance error structure with error code and message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// Error returns formatted error message.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// NewBinance returns Binance instance.
func NewBinance(service Service) Binance {
	return &binance{
		Service: service,
	}
}

// Ping tests connectivity.
func (b *binance) Ping() error {
	return b.Service.Ping()
}

// Time returns server time.
func (b *binance) Time() (time.Time, error) {
	return b.Service.Time()
}

// OrderBook represents Bids and Asks.
type OrderBook struct {
	LastUpdateID int `json:"lastUpdateId"`
	Bids         []*Order
	Asks         []*Order
}

type DepthEvent struct {
	WSEvent
	UpdateID int
	OrderBook
}

// Order represents single order information.
type Order struct {
	Price    float64
	Quantity float64
}

// OrderBookRequest represents OrderBook request data.
type OrderBookRequest struct {
	Symbol string
	Limit  int
}

// OrderBook returns list of orders.
func (b *binance) OrderBook(obr OrderBookRequest) (*OrderBook, error) {
	return b.Service.OrderBook(obr)
}

// AggTrade represents aggregated trade.
type AggTrade struct {
	ID             int
	Price          float64
	Quantity       float64
	FirstTradeID   int
	LastTradeID    int
	Timestamp      time.Time
	BuyerMaker     bool
	BestPriceMatch bool
}

type AggTradeEvent struct {
	WSEvent
	AggTrade
}

// AggTradesRequest represents AggTrades request data.
type AggTradesRequest struct {
	Symbol    string
	FromID    int64
	StartTime int64
	EndTime   int64
	Limit     int
}

// AggTrades returns compressed/aggregate list of trades.
func (b *binance) AggTrades(atr AggTradesRequest) ([]*AggTrade, error) {
	return b.Service.AggTrades(atr)
}

// KlinesRequest represents Klines request data.
type KlinesRequest struct {
	Symbol    string
	Interval  Interval
	Limit     int
	StartTime int64
	EndTime   int64
}

// Kline represents single Kline information.
type Kline struct {
	OpenTime                 time.Time
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                time.Time
	QuoteAssetVolume         float64
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

type KlineEvent struct {
	WSEvent
	Interval     Interval
	FirstTradeID int64
	LastTradeID  int64
	Final        bool
	Kline
}

// Klines returns klines/candlestick data.
func (b *binance) Klines(kr KlinesRequest) ([]*Kline, error) {
	return b.Service.Klines(kr)
}

// TickerRequest represents Ticker request data.
type TickerRequest struct {
	Symbol string
}

// Ticker24 represents data for 24hr ticker.
type Ticker24 struct {
	PriceChange        float64
	PriceChangePercent float64
	WeightedAvgPrice   float64
	PrevClosePrice     float64
	LastPrice          float64
	BidPrice           float64
	AskPrice           float64
	OpenPrice          float64
	HighPrice          float64
	LowPrice           float64
	Volume             float64
	OpenTime           time.Time
	CloseTime          time.Time
	FirstID            int
	LastID             int
	Count              int
}

// Ticker24 returns 24hr price change statistics.
func (b *binance) Ticker24(tr TickerRequest) (*Ticker24, error) {
	return b.Service.Ticker24(tr)
}

// PriceTicker represents ticker data for price.
type PriceTicker struct {
	Symbol string
	Price  float64
}

// TickerAllPrices returns ticker data for symbols.
func (b *binance) TickerAllPrices() ([]*PriceTicker, error) {
	return b.Service.TickerAllPrices()
}

// BookTicker represents book ticker data.
type BookTicker struct {
	Symbol   string
	BidPrice float64
	BidQty   float64
	AskPrice float64
	AskQty   float64
}

// TickerAllBooks returns tickers for all books.
func (b *binance) TickerAllBooks() ([]*BookTicker, error) {
	return b.Service.TickerAllBooks()
}

// NewOrderRequest represents NewOrder request data.
type NewOrderRequest struct {
	Symbol           string
	Side             OrderSide
	Type             OrderType
	TimeInForce      TimeInForce
	Quantity         float64
	Price            float64
	NewClientOrderID string
	StopPrice        float64
	IcebergQty       float64
	Timestamp        time.Time
}

// ProcessedOrder represents data from processed order.
type ProcessedOrder struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
	TransactTime  time.Time
}

// NewOrder places new order and returns ProcessedOrder.
func (b *binance) NewOrder(nor NewOrderRequest) (*ProcessedOrder, error) {
	return b.Service.NewOrder(nor)
}

// NewOrder places testing order.
func (b *binance) NewOrderTest(nor NewOrderRequest) error {
	return b.Service.NewOrderTest(nor)
}

// QueryOrderRequest represents QueryOrder request data.
type QueryOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

// ExecutedOrder represents data about executed order.
type ExecutedOrder struct {
	Symbol        string
	OrderID       int
	ClientOrderID string
	Price         float64
	OrigQty       float64
	ExecutedQty   float64
	Status        OrderStatus
	TimeInForce   TimeInForce
	Type          OrderType
	Side          OrderSide
	StopPrice     float64
	IcebergQty    float64
	Time          time.Time
}

// QueryOrder returns data about existing order.
func (b *binance) QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	return b.Service.QueryOrder(qor)
}

// CancelOrderRequest represents CancelOrder request data.
type CancelOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	NewClientOrderID  string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

// CanceledOrder represents data about canceled order.
type CanceledOrder struct {
	Symbol            string
	OrigClientOrderID string
	OrderID           int64
	ClientOrderID     string
}

// CancelOrder cancels order.
func (b *binance) CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	return b.Service.CancelOrder(cor)
}

// OpenOrdersRequest represents OpenOrders request data.
type OpenOrdersRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// OpenOrders returns list of open orders.
func (b *binance) OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	return b.Service.OpenOrders(oor)
}

// AllOrdersRequest represents AllOrders request data.
type AllOrdersRequest struct {
	Symbol     string
	OrderID    int64
	Limit      int
	RecvWindow time.Duration
	Timestamp  time.Time
}

// AllOrders returns list of all previous orders.
func (b *binance) AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	return b.Service.AllOrders(aor)
}

// AccountRequest represents Account request data.
type AccountRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Account represents user's account information.
type Account struct {
	MakerCommision  int64
	TakerCommision  int64
	BuyerCommision  int64
	SellerCommision int64
	CanTrade        bool
	CanWithdraw     bool
	CanDeposit      bool
	Balances        []*Balance
}

type AccountEvent struct {
	WSEvent
	Account
}

// Balance groups balance-related information.
type Balance struct {
	Asset  string
	Free   float64
	Locked float64
}

// Account returns account data.
func (b *binance) Account(ar AccountRequest) (*Account, error) {
	return b.Service.Account(ar)
}

// MyTradesRequest represents MyTrades request data.
type MyTradesRequest struct {
	Symbol     string
	Limit      int
	FromID     int64
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Trade represents data about trade.
type Trade struct {
	ID              int64
	Price           float64
	Qty             float64
	Commission      float64
	CommissionAsset string
	Time            time.Time
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}

// MyTrades list user's trades.
func (b *binance) MyTrades(mtr MyTradesRequest) ([]*Trade, error) {
	return b.Service.MyTrades(mtr)
}

// WithdrawRequest represents Withdraw request data.
type WithdrawRequest struct {
	Asset      string
	Address    string
	Amount     float64
	Name       string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// WithdrawResult represents Withdraw result.
type WithdrawResult struct {
	Success bool
	Msg     string
}

// Withdraw executes withdrawal.
func (b *binance) Withdraw(wr WithdrawRequest) (*WithdrawResult, error) {
	return b.Service.Withdraw(wr)
}

// HistoryRequest represents history-related calls request data.
type HistoryRequest struct {
	Asset      string
	Status     *int
	StartTime  time.Time
	EndTime    time.Time
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Deposit represents Deposit data.
type Deposit struct {
	InsertTime time.Time
	Amount     float64
	Asset      string
	Status     int
}

// DepositHistory lists deposit data.
func (b *binance) DepositHistory(hr HistoryRequest) ([]*Deposit, error) {
	return b.Service.DepositHistory(hr)
}

// Withdrawal represents withdrawal data.
type Withdrawal struct {
	Amount    float64
	Address   string
	TxID      string
	Asset     string
	ApplyTime time.Time
	Status    int
}

// WithdrawHistory lists withdraw data.
func (b *binance) WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error) {
	return b.Service.WithdrawHistory(hr)
}


type RecentTradesRequest struct {
	Symbol  string
	Limit   int
}

type RecentTrades struct {
	ID              int64
	Price           float64
	Qty             float64
	Time            time.Time
	IsBuyerMaker    bool
	IsBestMatch     bool
}

// Recent trades list.
func (b *binance) RecentTrades(rt RecentTradesRequest) ([]*RecentTrades, error) {
	return b.Service.RecentTrades(rt)
}


// Stream represents stream information.
//
// Read web docs to get more information about using streams.
type Stream struct {
	ListenKey string
}

// StartUserDataStream starts stream and returns Stream with ListenKey.
func (b *binance) StartUserDataStream() (*Stream, error) {
	return b.Service.StartUserDataStream()
}

// KeepAliveUserDataStream prolongs stream livespan.
func (b *binance) KeepAliveUserDataStream(s *Stream) error {
	return b.Service.KeepAliveUserDataStream(s)
}

// CloseUserDataStream closes opened stream.
func (b *binance) CloseUserDataStream(s *Stream) error {
	return b.Service.CloseUserDataStream(s)
}

type WSEvent struct {
	Type   string
	Time   time.Time
	Symbol string
}

type DepthWebsocketRequest struct {
	Symbol string
}

func (b *binance) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	return b.Service.DepthWebsocket(dwr)
}

type KlineWebsocketRequest struct {
	Symbol   string
	Interval Interval
}

func (b *binance) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	return b.Service.KlineWebsocket(kwr)
}

type TradeWebsocketRequest struct {
	Symbol string
}

func (b *binance) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	return b.Service.TradeWebsocket(twr)
}

type UserDataWebsocketRequest struct {
	ListenKey string
}

func (b *binance) UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	return b.Service.UserDataWebsocket(udwr)
}
