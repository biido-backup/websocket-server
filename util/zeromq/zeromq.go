package zeromq

import (
	"context"
	"encoding/json"
	"github.com/go-zeromq/zmq4"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"websocket-server/cache"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/logger"
	"websocket-server/util/websocket"
)

var log = logger.CreateLog("zeromq")

var HasConnectedZeroMQOrderBook = false
var HasConnectedZeroMQLast24h = false
var HasConnectedZeroMQTradingHistory = false
var HasConnectedZeroMQOpenOrder = false
var HasConnectedZeroMQOrderHistory = false

func Listen(tradingRateList []daos.Rate)  {
	done := make(chan bool)

	c := cron.New()
	c.AddFunc("0-59/5 * * * * *", func() {
		ListenAll(tradingRateList)
	})
	c.Start()

	<- done
}

func ListenAll(tradingRateList []daos.Rate) {
	if !HasConnectedZeroMQOrderBook {
		log.Debug("Reconnect to zeromq order book")
		for _, tradingRate := range(tradingRateList){
			go ListenOrderBook(tradingRate.StringDash())
		}
	}

	if !HasConnectedZeroMQLast24h {
		log.Debug("Reconnect to zeromq last 24h")
		for _, tradingRate := range(tradingRateList){
			go ListenLast24h(tradingRate.StringDash())
		}
	}

	if !HasConnectedZeroMQTradingHistory {
		log.Debug("Reconnect to zeromq trading history")
		for _, tradingRate := range(tradingRateList){
			go ListenTradingHistory(tradingRate.StringDash())
		}
	}

	if !HasConnectedZeroMQOpenOrder {
		log.Debug("Reconnect to zeromq open order")
		for _, tradingRate := range(tradingRateList){
			go ListenOpenOrder(tradingRate.StringDash())
		}
	}

	if !HasConnectedZeroMQOrderHistory {
		log.Debug("Reconnect to zeromq order history")
		for _, tradingRate := range(tradingRateList){
			go ListenOrderHistory(tradingRate.StringDash())
		}
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

func ListenOrderBook(topic string){
	publisher := viper.GetString("zeromq.publisher.matching-engine")
	zmqKey := topic + viper.GetString("zeromq.key.suffix.orderbook")

	(&daos.MyClients).SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQOrderBook = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQOrderBook = false
		return
	}

	HasConnectedZeroMQOrderBook = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQOrderBook = false
			break
		}
		//
		var orderbook trading.OrderBookZeroMQ
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		orderbook = trading.OrderBookFromJSONZeroMQ(msg.Frames[1])

		trdOrderbook := trading.Orderbook{trdconst.ORDERBOOK, orderbook}
		trdOrderbookJson, err := json.Marshal(trdOrderbook)
		if err!=nil{
			log.Error("error when unmarshal orderbook", err)
		}

		//log.Debug(orderbook)
		cache.SetCacheByTopicAndType(topic, trdconst.ORDERBOOK, trdOrderbook)
		websocket.BroadcastMessageWithTopic(topic, string(trdOrderbookJson))
	}
}


//func ListenLast24h(publisher string, topic string, zmqKey string, clients *daos.Clients){
func ListenLast24h(topic string){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := topic + viper.GetString("zeromq.key.suffix.last24h")

	(&daos.MyClients).SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQLast24h = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQLast24h = false
		return
	}

	HasConnectedZeroMQLast24h = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQLast24h = false
			break
		}
		//

		var last24h trading.Last24h
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &last24h)
		if err != nil {
			log.Error("error when unmarshal last24h", err)
		}

		trdLast24h := trading.TradingLast24h{trdconst.LAST24H, topic, last24h}
		trdLast24hJson, err := json.Marshal(trdLast24h)
		if err!=nil{
			log.Error("error when marshal last24h", err)
		}

		//log.Debug(trdLast24h)
		cache.SetCacheByTopicAndType(topic, trdconst.LAST24H, trdLast24h)
		websocket.BroadcastMessageToAll(string(trdLast24hJson))
	}
}


func ListenTradingHistory(topic string){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := topic + viper.GetString("zeromq.key.suffix.tradinghistory")

	(&daos.MyClients).SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQTradingHistory = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQTradingHistory = false
		return
	}

	HasConnectedZeroMQTradingHistory = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQTradingHistory = false
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

		trdListHistory := trading.TradingListHistory{trdconst.TRADINGHISTORY, listHistory}
		trdListHistoryJson, err := json.Marshal(trdListHistory)
		if err!=nil{
			log.Error("error when marshal listHistory", err)
		}

		//log.Debug(trdListHistory)
		cache.SetCacheByTopicAndType(topic, trdconst.TRADINGHISTORY, trdListHistory)
		websocket.BroadcastMessageWithTopic(topic, string(trdListHistoryJson))
	}
}

func ListenOpenOrder(topic string){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := topic + viper.GetString("zeromq.key.suffix.openorder")

	(&daos.MyClients).SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQOpenOrder = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQOpenOrder = false
		return
	}

	HasConnectedZeroMQOpenOrder = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQOpenOrder = false
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

		trdListOpenOrder := trading.OpenOrders{Type:trdconst.OPENORDER, Payload:listOpenOrder.OpenOrders}
		trdListOpenOrderJson, err := json.Marshal(trdListOpenOrder)
		if err!=nil{
			log.Error("error when marshal listOpenOrder", err)
		}

		username := listOpenOrder.Username
		//log.Println("username : ",username)
		//log.Debug(trdListOpenOrder)
		websocket.SendMessageToUser(topic, username, string(trdListOpenOrderJson))
	}
}

func ListenOrderHistory(topic string){
	publisher := viper.GetString("zeromq.publisher.trading-broker")
	zmqKey := topic + viper.GetString("zeromq.key.suffix.orderhistory")

	const SIZE = 5

	(&daos.MyClients).SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQOrderHistory = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQOrderHistory = false
		return
	}

	HasConnectedZeroMQOrderHistory = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQOrderHistory = false
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

		trdListOrderHistory := trading.OrderHistories{Type:trdconst.ORDERHISTORY, Size:listOrderHistory.Size, Payload:listOrderHistory.OrderHistories}
		trdListOrderHistoryJson, err := json.Marshal(trdListOrderHistory)
		if err!=nil{
			log.Error("error when marshal listOrderHistory", err)
		}

		username := listOrderHistory.Username
		//log.Println("username : ",username)
		//log.Debug(trdListOrderHistory)
		websocket.SendMessageToUser(topic, username, string(trdListOrderHistoryJson))
	}
}