package zeromq

import (
	"context"
	"encoding/json"
	"github.com/go-zeromq/zmq4"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"websocket-server/cache"
	"websocket-server/const/trd"
	"websocket-server/daos/trading"
	"websocket-server/util/logger"
	"websocket-server/util/websocket"
)

var log = logger.CreateLog("zeromq")

var ConnectedToZeroMQOrderBook = false
var ConnectedToZeroMQLast24H = false
var ConnectedToZeroMQTradingHistory = false
var ConnectedToZeroMQOpenOrder = false
var ConnectedToZeroMQOrderHistory = false


func Listen()  {
	done := make(chan bool)

	c := cron.New()
	c.AddFunc("0-59/5 * * * * *", func() {
		ListenAll()
	})
	c.Start()

	<- done
}

func ListenAll() {
	if !ConnectedToZeroMQOrderBook {
		log.Info("Reconnect to zeromq order book")
		//for _, tradingRate := range(tradingRateList){
		//	go ListenOrderBook(tradingRate.StringDash())
		//}
		go ListenOrderBook()
	}

	if !ConnectedToZeroMQLast24H {
		log.Info("Reconnect to zeromq last 24h")
		//for _, tradingRate := range(tradingRateList){
		//	go ListenLast24h(tradingRate.StringDash())
		//}
		go ListenLast24h()
	}

	if !ConnectedToZeroMQTradingHistory {
		log.Info("Reconnect to zeromq trading history")
		//for _, tradingRate := range(tradingRateList){
		//	go ListenTradingHistory(tradingRate.StringDash())
		//}
		go ListenTradingHistory()
	}

	if !ConnectedToZeroMQOpenOrder {
		log.Debug("Reconnect to zeromq open order")
		//for _, tradingRate := range(tradingRateList){
		//	go ListenOpenOrder(tradingRate.StringDash())
		//}
		go ListenOpenOrder()
	}

	if !ConnectedToZeroMQOrderHistory {
		log.Debug("Reconnect to zeromq order history")
		//for _, tradingRate := range(tradingRateList){
		//	go ListenOrderHistory(tradingRate.StringDash())
		//}
		go ListenOrderHistory()
	}
}

func ListenTest(publisher string, zmqKey string){
	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher", err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)

	if err != nil {
		log.Error("could not subscribe publisher", err)
		return
	}

	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)
			continue
		}
		//
		var orderbook trading.OrderBookZeroMQ
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		orderbook = trading.OrderBookFromJSONZeroMQ(msg.Frames[1])

		log.Debug(orderbook)

	}
}

func ListenOrderBook(){
	publisher := viper.GetString("zeromq.publisher.matching-engine")
	zmqKey := viper.GetString("zeromq.key.suffix.orderbook")

	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		ConnectedToZeroMQOrderBook = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
		ConnectedToZeroMQOrderBook = false
		return
	}

	ConnectedToZeroMQOrderBook = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)
			sub.Close()
			ConnectedToZeroMQOrderBook = false
			break
		}
		//
		var orderbook trading.OrderBookZeroMQ
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		orderbook = trading.OrderBookFromJSONZeroMQ(msg.Frames[1])
		topic := orderbook.Rate

		trdOrderbook := trading.Orderbook{trdconst.ORDERBOOK, orderbook}
		trdOrderbookJson, err := json.Marshal(trdOrderbook)
		if err!=nil{
			log.Error("error when marshal orderbook", err)
		}

		//log.Debug(orderbook)
		cache.SetCacheByTopicAndType(topic, trdconst.ORDERBOOK, trdOrderbook)
		websocket.BroadcastMessageWithTopic(orderbook.Rate, string(trdOrderbookJson))
	}
}


//func ListenLast24h(publisher string, topic string, zmqKey string, clients *daos.Clients){
func ListenLast24h(){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := viper.GetString("zeromq.key.suffix.last24h")

	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		ConnectedToZeroMQLast24H = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		ConnectedToZeroMQLast24H = false
		return
	}

	ConnectedToZeroMQLast24H = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			ConnectedToZeroMQLast24H = false
			break
		}
		//

		var last24hTrx trading.Last24hTrx
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &last24hTrx)
		if err != nil {
			log.Error("error when unmarshal last24hTrx", err)
		}

		topic := last24hTrx.Rate

		trdLast24h := trading.TradingLast24h{trdconst.LAST24H, topic, last24hTrx.Last24h}
		trdLast24hJson, err := json.Marshal(trdLast24h)
		if err!=nil{
			log.Error("error when marshal last24hTrx", err)
		}

		log.Debug("ListenLast24h", trdLast24h.Payload)
		cache.SetCacheByTopicAndType(topic, trdconst.LAST24H, trdLast24h)
		websocket.BroadcastMessageToAll(string(trdLast24hJson))
	}
}


func ListenTradingHistory(){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := viper.GetString("zeromq.key.suffix.tradinghistory")

	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		ConnectedToZeroMQTradingHistory = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		ConnectedToZeroMQTradingHistory = false
		return
	}

	ConnectedToZeroMQTradingHistory = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			ConnectedToZeroMQTradingHistory = false
			break
		}
		//

		var listHistory trading.ListHistory
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &listHistory)
		if err != nil {
			log.Error("error when unmarshal listHistory", err)
		}

		topic := listHistory.Rate

		trdListHistory := trading.TradingListHistory{trdconst.TRADINGHISTORY, listHistory.Histories}
		trdListHistoryJson, err := json.Marshal(trdListHistory)
		if err!=nil{
			log.Error("error when marshal listHistory", err)
		}

		log.Debug("ListenTradingHistory",trdListHistory.Payload[0])
		cache.SetCacheByTopicAndType(topic, trdconst.TRADINGHISTORY, trdListHistory)
		websocket.BroadcastMessageWithTopic(topic, string(trdListHistoryJson))
	}
}

func ListenOpenOrder(){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := viper.GetString("zeromq.key.suffix.openorder")

	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		ConnectedToZeroMQOpenOrder = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		ConnectedToZeroMQOpenOrder = false
		return
	}

	ConnectedToZeroMQOpenOrder = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			ConnectedToZeroMQOpenOrder = false
			break
		}
		//

		var listOpenOrder trading.ListOpenOrder
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &listOpenOrder)
		if err != nil {
			log.Error("error when unmarshal listOpenOrder", err)
		}

		topic := listOpenOrder.Rate

		trdListOpenOrder := trading.OpenOrders{Type:trdconst.OPENORDER, Payload:listOpenOrder.OpenOrders}
		trdListOpenOrderJson, err := json.Marshal(trdListOpenOrder)
		if err!=nil{
			log.Error("error when marshal listOpenOrder", err)
		}

		username := listOpenOrder.Username
		//log.Println("username : ",username)
		log.Debug("ListenOpenOrder", trdListOpenOrder.Payload[0])
		websocket.SendMessageToUser(topic, username, string(trdListOpenOrderJson))
	}
}

func ListenOrderHistory(){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := viper.GetString("zeromq.key.suffix.orderhistory")

	const SIZE = 5

	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		ConnectedToZeroMQOrderHistory = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		ConnectedToZeroMQOrderHistory = false
		return
	}

	ConnectedToZeroMQOrderHistory = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			ConnectedToZeroMQOrderHistory = false
			break
		}
		//

		var listOrderHistory trading.ListOrderHistory
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &listOrderHistory)
		if err != nil {
			log.Error("error when unmarshal listOrderHistory", err)
		}

		topic := listOrderHistory.Rate

		trdListOrderHistory := trading.OrderHistories{Type:trdconst.ORDERHISTORY, Size:listOrderHistory.Size, Payload:listOrderHistory.OrderHistories}
		trdListOrderHistoryJson, err := json.Marshal(trdListOrderHistory)
		if err!=nil{
			log.Error("error when marshal listOrderHistory", err)
		}

		username := listOrderHistory.Username
		//log.Println("username : ",username)
		log.Debug("ListenOrderHistory", trdListOrderHistory.Payload[0])
		websocket.SendMessageToUser(topic, username, string(trdListOrderHistoryJson))
	}
}