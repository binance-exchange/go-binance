# Binance API

To read full documentation, specs and find out which request params are required/optional, please visit the official
[documentation](https://www.binance.com/restapipub.html) page.

## Getting started

```go
var logger log.Logger
logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

hmacSigner := &binance.HmacSigner{
    Key: []byte("API secret"),
}
ctx, _ := context.WithCancel(context.Background())
// use second return value for cancelling request when shutting down the app

binanceService := binance.NewAPIService(
    "https://www.binance.com",
    "API key",
    hmacSigner,
    logger,
    ctx,
)
b := binance.NewBinance(binanceService)
```

## Examples

Following provides list of main usages of library. See `example` package for testing application with more examples.

Each call has its own *Request* structure with data that can be provided. The library is not responsible for validating
the input and if non-zero value is used, the param is sent to the API server.

In case of an standard error, instance of `binance.Error` is returned with additional info.

### NewOrder

```go
newOrder, err := b.NewOrder(binance.NewOrderRequest{
    Symbol:      "BNBETH",
    Quantity:    1,
    Price:       999,
    Side:        binance.SideSell,
    TimeInForce: binance.GTC,
    Type:        binance.TypeLimit,
    Timestamp:   time.Now(),
})
if err != nil {
    panic(err)
}
fmt.Println(newOrder)
```

### CancelOrder

```go
canceledOrder, err := b.CancelOrder(binance.CancelOrderRequest{
    Symbol:    "BNBETH",
    OrderID:   newOrder.OrderID,
    Timestamp: time.Now(),
})
if err != nil {
    panic(err)
}
fmt.Printf("%#v\n", canceledOrder)
```

### Klines

```go
kl, err := b.Klines(binance.KlinesRequest{
    Symbol:   "BNBETH",
    Interval: binance.Hour,
})
if err != nil {
    panic(err)
}
fmt.Printf("%#v\n", kl)
```
    
### Trade Websocket

```go
interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt)

kech, done, err := b.TradeWebsocket(binance.TradeWebsocketRequest{
    Symbol: "ETHBTC",
})
if err != nil {
    panic(err)
}
go func() {
    for {
        select {
        case ke := <-kech:
            fmt.Printf("%#v\n", ke)
        case <-done:
            break
        }
    }
}()

fmt.Println("waiting for interrupt")
<-interrupt
fmt.Println("canceling context")
cancelCtx()
fmt.Println("waiting for signal")
<-done
fmt.Println("exit")
return
```

## Known issues

* Websocket error handling is not perfect and occasionally attempts to read from closed connection.