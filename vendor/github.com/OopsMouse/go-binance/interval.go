package binance

// Interval represents interval enum.
type Interval string

var (
	Minute         = Interval("1m")
	ThreeMinutes   = Interval("3m")
	FiveMinutes    = Interval("5m")
	FifteenMinutes = Interval("15m")
	ThirtyMinutes  = Interval("30m")
	Hour           = Interval("1h")
	TwoHours       = Interval("2h")
	FourHours      = Interval("4h")
	SixHours       = Interval("6h")
	EightHours     = Interval("8h")
	TwelveHours    = Interval("12h")
	Day            = Interval("1d")
	ThreeDays      = Interval("3d")
	Week           = Interval("1w")
	Month          = Interval("1M")
)

// Time in force to indicate how long an order will remain active before it is executed or expires.
type TimeInForce string

var (
	// GTC (Good-Til-Canceled) orders are effective until they are executed or canceled.
	GTC = TimeInForce("GTC")
	// IOC (Immediate or Cancel) orders fills all or part of an order immediately and cancels the remaining part of the order.
	IOC = TimeInForce("IOC")
	// FOK (Fill or Kill) orders fills all in its entirety, otherwise, the entire order will be cancelled.
	FOK = TimeInForce("FOK")
)
