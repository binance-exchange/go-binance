package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
)

func (as *apiService) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth", strings.ToLower(dwr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	dech := make(chan *DepthEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawDepth := struct {
					Type          string          `json:"e"`
					Time          float64         `json:"E"`
					Symbol        string          `json:"s"`
					UpdateID      int             `json:"u"`
					BidDepthDelta [][]interface{} `json:"b"`
					AskDepthDelta [][]interface{} `json:"a"`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := timeFromUnixTimestampFloat(rawDepth.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				de := &DepthEvent{
					WSEvent: WSEvent{
						Type:   rawDepth.Type,
						Time:   t,
						Symbol: rawDepth.Symbol,
					},
					UpdateID: rawDepth.UpdateID,
				}
				for _, b := range rawDepth.BidDepthDelta {
					p, err := floatFromString(b[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}
					q, err := floatFromString(b[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}
					de.Bids = append(de.Bids, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				dech <- de
			}
		}
	}()

	go as.exitHandler(c, done)
	return dech, done, nil
}

func (as *apiService) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", strings.ToLower(kwr.Symbol), string(kwr.Interval))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	kech := make(chan *KlineEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawKline := struct {
					Type     string  `json:"e"`
					Time     float64 `json:"E"`
					Symbol   string  `json:"S"`
					OpenTime float64 `json:"t"`
					Kline    struct {
						Interval                 string  `json:"i"`
						FirstTradeID             int64   `json:"f"`
						LastTradeID              int64   `json:"L"`
						Final                    bool    `json:"x"`
						OpenTime                 float64 `json:"t"`
						CloseTime                float64 `json:"T"`
						Open                     string  `json:"o"`
						High                     string  `json:"h"`
						Low                      string  `json:"l"`
						Close                    string  `json:"c"`
						Volume                   string  `json:"v"`
						NumberOfTrades           int     `json:"n"`
						QuoteAssetVolume         string  `json:"q"`
						TakerBuyBaseAssetVolume  string  `json:"V"`
						TakerBuyQuoteAssetVolume string  `json:"Q"`
					} `json:"k"`
				}{}
				if err := json.Unmarshal(message, &rawKline); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := timeFromUnixTimestampFloat(rawKline.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Time)
					return
				}
				ot, err := timeFromUnixTimestampFloat(rawKline.Kline.OpenTime)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.OpenTime)
					return
				}
				ct, err := timeFromUnixTimestampFloat(rawKline.Kline.CloseTime)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.CloseTime)
					return
				}
				open, err := floatFromString(rawKline.Kline.Open)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Open)
					return
				}
				cls, err := floatFromString(rawKline.Kline.Close)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Close)
					return
				}
				high, err := floatFromString(rawKline.Kline.High)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.High)
					return
				}
				low, err := floatFromString(rawKline.Kline.Low)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Low)
					return
				}
				vol, err := floatFromString(rawKline.Kline.Volume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Volume)
					return
				}
				qav, err := floatFromString(rawKline.Kline.QuoteAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", (rawKline.Kline.QuoteAssetVolume))
					return
				}
				tbbav, err := floatFromString(rawKline.Kline.TakerBuyBaseAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.TakerBuyBaseAssetVolume)
					return
				}
				tbqav, err := floatFromString(rawKline.Kline.TakerBuyQuoteAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.TakerBuyQuoteAssetVolume)
					return
				}

				ke := &KlineEvent{
					WSEvent: WSEvent{
						Type:   rawKline.Type,
						Time:   t,
						Symbol: rawKline.Symbol,
					},
					Interval:     Interval(rawKline.Kline.Interval),
					FirstTradeID: rawKline.Kline.FirstTradeID,
					LastTradeID:  rawKline.Kline.LastTradeID,
					Final:        rawKline.Kline.Final,
					Kline: Kline{
						OpenTime:                 ot,
						CloseTime:                ct,
						Open:                     open,
						Close:                    cls,
						High:                     high,
						Low:                      low,
						Volume:                   vol,
						NumberOfTrades:           rawKline.Kline.NumberOfTrades,
						QuoteAssetVolume:         qav,
						TakerBuyBaseAssetVolume:  tbbav,
						TakerBuyQuoteAssetVolume: tbqav,
					},
				}
				kech <- ke
			}
		}
	}()

	go as.exitHandler(c, done)
	return kech, done, nil
}

func (as *apiService) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", strings.ToLower(twr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aggtech := make(chan *AggTradeEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawAggTrade := struct {
					Type         string  `json:"e"`
					Time         float64 `json:"E"`
					Symbol       string  `json:"s"`
					TradeID      int     `json:"a"`
					Price        string  `json:"p"`
					Quantity     string  `json:"q"`
					FirstTradeID int     `json:"f"`
					LastTradeID  int     `json:"l"`
					Timestamp    float64 `json:"T"`
					IsMaker      bool    `json:"m"`
				}{}
				if err := json.Unmarshal(message, &rawAggTrade); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := timeFromUnixTimestampFloat(rawAggTrade.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAggTrade.Time)
					return
				}

				price, err := floatFromString(rawAggTrade.Price)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAggTrade.Price)
					return
				}
				qty, err := floatFromString(rawAggTrade.Quantity)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAggTrade.Quantity)
					return
				}
				ts, err := timeFromUnixTimestampFloat(rawAggTrade.Timestamp)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAggTrade.Timestamp)
					return
				}

				ae := &AggTradeEvent{
					WSEvent: WSEvent{
						Type:   rawAggTrade.Type,
						Time:   t,
						Symbol: rawAggTrade.Symbol,
					},
					AggTrade: AggTrade{
						ID:           rawAggTrade.TradeID,
						Price:        price,
						Quantity:     qty,
						FirstTradeID: rawAggTrade.FirstTradeID,
						LastTradeID:  rawAggTrade.LastTradeID,
						Timestamp:    ts,
						BuyerMaker:   rawAggTrade.IsMaker,
					},
				}
				aggtech <- ae
			}
		}
	}()

	go as.exitHandler(c, done)
	return aggtech, done, nil
}

func (as *apiService) UserDataWebsocket(urwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s", urwr.ListenKey)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aech := make(chan *AccountEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}

				// get payload type
				payload := struct {
					Type string `json:"e"`
				}{}
				if err := json.Unmarshal(message, &payload); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}

				level.Error(as.Logger).Log("payload", payload, "type", payload.Type)

				// deal with different payload types
				if payload.Type == "outboundAccountInfo" {
					rawAccount := struct {
						Type            string  `json:"e"`
						Time            float64 `json:"E"`
						OpenTime        float64 `json:"t"`
						MakerCommision  int64   `json:"m"`
						TakerCommision  int64   `json:"t"`
						BuyerCommision  int64   `json:"b"`
						SellerCommision int64   `json:"s"`
						CanTrade        bool    `json:"T"`
						CanWithdraw     bool    `json:"W"`
						CanDeposit      bool    `json:"D"`
						Balances        []struct {
							Asset            string `json:"a"`
							AvailableBalance string `json:"f"`
							Locked           string `json:"l"`
						} `json:"B"`
					}{}

					if err := json.Unmarshal(message, &rawAccount); err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}

					t, err := timeFromUnixTimestampFloat(rawAccount.Time)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAccount.Time)
						return
					}

					ae := &AccountEvent{
						WSEvent: WSEvent{
							Type: rawAccount.Type,
							Time: t,
						},
						Account: Account{
							MakerCommision:  rawAccount.MakerCommision,
							TakerCommision:  rawAccount.TakerCommision,
							BuyerCommision:  rawAccount.BuyerCommision,
							SellerCommision: rawAccount.SellerCommision,
							CanTrade:        rawAccount.CanTrade,
							CanWithdraw:     rawAccount.CanWithdraw,
							CanDeposit:      rawAccount.CanDeposit,
						},
						ExecutedOrder: ExecutedOrder{},
					}
					for _, b := range rawAccount.Balances {
						free, err := floatFromString(b.AvailableBalance)
						if err != nil {
							level.Error(as.Logger).Log("wsUnmarshal", err, "body", b.AvailableBalance)
							return
						}
						locked, err := floatFromString(b.Locked)
						if err != nil {
							level.Error(as.Logger).Log("wsUnmarshal", err, "body", b.Locked)
							return
						}
						ae.Balances = append(ae.Balances, &Balance{
							Asset:  b.Asset,
							Free:   free,
							Locked: locked,
						})
					}

					// submit data to channel
					aech <- ae
				} else if payload.Type == "executionReport" {
					rawOrderUpdate := struct {
						Type                     string  `json:"e"`
						Time                     float64 `json:"E"`
						Symbol                   string  `json:"s"`
						ClientOrderID            string  `json:"c"`
						Side                     string  `json:"S"`
						OrderType                string  `json:"o"`
						TimeInForce              string  `json:"f"`
						Quantity                 string  `json:"q"`
						Price                    string  `json:"p"`
						StopPrice                string  `json:"P"`
						IcebergQty               string  `json:"F"`
						OriginalClientOrderID    string  `json:"C"`
						ExecutionType            string  `json:"x"`
						Status                   string  `json:"X"`
						OrderRejectReason        string  `json:"r"`
						OrderID                  int     `json:"i"`
						LastExecutedQuantity     string  `json:"l"`
						CumulativeFilledQuantity string  `json:"z"`
						LastExecutedPrice        string  `json:"L"`
						CommissionAmount         string  `json:"n"`
						TransactionTime          float64 `json:"T"`
						TradeID                  float64 `json:"t"`
					}{}

					if err := json.Unmarshal(message, &rawOrderUpdate); err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}

					t, err := timeFromUnixTimestampFloat(rawOrderUpdate.Time)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawOrderUpdate.Time)
						return
					}

					price, err := strconv.ParseFloat(rawOrderUpdate.Price, 64)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", "cannot parse Price")
						return
					}
					origQty, err := strconv.ParseFloat(rawOrderUpdate.Quantity, 64)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", "cannot parse origQty")
						return
					}
					execQty, err := strconv.ParseFloat(rawOrderUpdate.LastExecutedQuantity, 64)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", "cannot parse execQty")
						return
					}
					stopPrice, err := strconv.ParseFloat(rawOrderUpdate.StopPrice, 64)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", "cannot parse stopPrice")
						return
					}
					icebergQty, err := strconv.ParseFloat(rawOrderUpdate.IcebergQty, 64)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", "cannot parse icebergQty")
						return
					}

					ae := &AccountEvent{
						WSEvent: WSEvent{
							Type: rawOrderUpdate.Type,
							Time: t,
						},
						Account: Account{},
						ExecutedOrder: ExecutedOrder{
							Symbol:        rawOrderUpdate.Symbol,
							OrderID:       rawOrderUpdate.OrderID,
							ClientOrderID: rawOrderUpdate.ClientOrderID,
							Price:         price,
							OrigQty:       origQty,
							ExecutedQty:   execQty,
							Status:        OrderStatus(rawOrderUpdate.Status),
							TimeInForce:   TimeInForce(rawOrderUpdate.TimeInForce),
							Type:          OrderType(rawOrderUpdate.Type),
							Side:          OrderSide(rawOrderUpdate.Side),
							StopPrice:     stopPrice,
							IcebergQty:    icebergQty,
							Time:          t,
						},
					}

					// submit data to channel
					aech <- ae
				}
			}
		}
	}()

	go as.exitHandler(c, done)
	return aech, done, nil
}

func (as *apiService) exitHandler(c *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				level.Error(as.Logger).Log("wsWrite", err)
				return
			}
		case <-as.Ctx.Done():
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			level.Info(as.Logger).Log("closing connection")
			return
		}
	}
}
