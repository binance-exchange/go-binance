package binance

// OrderStatus represents order status enum.
type OrderStatus string

// OrderType represents order type enum.
type OrderType string

// OrderSide represents order side enum.
type OrderSide string

var (
	StatusNew             = OrderStatus("NEW")
	StatusPartiallyFilled = OrderStatus("PARTIALLY_FILLED")
	StatusFilled          = OrderStatus("FILLED")
	StatusCancelled       = OrderStatus("CANCELED")
	StatusPendingCancel   = OrderStatus("PENDING_CANCEL")
	StatusRejected        = OrderStatus("REJECTED")
	StatusExpired         = OrderStatus("EXPIRED")

	// LIMIT	timeInForce, quantity, price
	TypeLimit = OrderType("LIMIT")
	// MARKET	quantity
	TypeMarket = OrderType("MARKET")
	// STOP_LOSS	quantity, stopPrice
	TypeStopLoss = OrderType("STOP_LOSS")
	// STOP_LOSS_LIMIT	timeInForce, quantity, price, stopPrice
	TypeStopLossLimit = OrderType("STOP_LOSS_LIMIT")
	// TAKE_PROFIT	quantity, stopPrice
	TypeTakeProfit = OrderType("TAKE_PROFIT")
	// TAKE_PROFIT_LIMIT	timeInForce, quantity, price, stopPrice
	TypeTakeProfitLimit = OrderType("TAKE_PROFIT_LIMIT")
	// LIMIT_MAKER	quantity, price
	TypeLimitMaker = OrderType("LIMIT_MAKER")

	SideBuy  = OrderSide("BUY")
	SideSell = OrderSide("SELL")
)

func (o OrderStatus) String() string {
	return string(o)
}
