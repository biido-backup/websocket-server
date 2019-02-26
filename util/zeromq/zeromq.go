package zeromq

import (
	"context"
	"encoding/json"
	"github.com/go-zeromq/zmq4"
	"github.com/spf13/viper"
	"time"
	"websocket-server/cache"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/logger"
	"websocket-server/util/websocket"
)

var log = logger.CreateLog("zeromq")

//func Listen(topic string, clients *daos.Clients){
func Listen(topic string){
	var matchingEngineAddr = viper.GetString("zeromq.publisher.matching-engine")
	var tradingBrokerAddr = viper.GetString("zeromq.publisher.trading-broker")

	var suffixOrderBook = viper.GetString("zeromq.key.suffix.orderbook")
	var suffixLast24h = viper.GetString("zeromq.key.suffix.last24h")
	var suffixTradingHistory = viper.GetString("zeromq.key.suffix.tradinghistory")
	var suffixOpenOrder = viper.GetString("zeromq.key.suffix.openorder")
	var suffixOrderHistory = viper.GetString("zeromq.key.suffix.orderhistory")
	//

	//log.Println(topic)

	//go ListenTest("tcp://localhost:5563", "BTC-IDR:ORDER_BOOK")
	//go ListenTest("tcp://localhost:5564", "BTC-IDR:LAST24H_TRANSACTION")
	//go ListenTest("tcp://localhost:5564", "BTC-IDR:TRADING_HISTORY")

	go ListenOrderBook(matchingEngineAddr, topic, topic+suffixOrderBook, &daos.MyClients)
	time.Sleep(time.Millisecond)
	go ListenLast24h(tradingBrokerAddr, topic, topic+suffixLast24h, &daos.MyClients)
	time.Sleep(time.Millisecond)
	go ListenTradingHistory(tradingBrokerAddr, topic, topic+suffixTradingHistory, &daos.MyClients)
	time.Sleep(time.Millisecond)
	go ListenOpenOrder(tradingBrokerAddr, topic, topic+suffixOpenOrder, &daos.MyClients)
	time.Sleep(time.Millisecond)
	go ListenOrderHistory(tradingBrokerAddr, topic, topic+suffixOrderHistory, &daos.MyClients)
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

		log.Println(orderbook)

	}
}

func ListenOrderBook(publisher string, topic string, zmqKey string, clients *daos.Clients){

	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
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

		trdOrderbook := trading.Orderbook{trdconst.ORDERBOOK, orderbook}
		trdOrderbookJson, err := json.Marshal(trdOrderbook)
		if err!=nil{
			log.Error("error when unmarshal orderbook", err)
		}

		log.Println(orderbook)
		cache.SetCacheByTopicAndType(topic, trdconst.ORDERBOOK, trdOrderbook)

		websocket.BroadcastMessage(topic, string(trdOrderbookJson))

	}

}


func ListenLast24h(publisher string, topic string, zmqKey string, clients *daos.Clients){

	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
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

		var last24h trading.Last24h
		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		err = json.Unmarshal(msg.Frames[1], &last24h)
		if err != nil {
			log.Error("error when unmarshal last24h", err)
		}

		trdLast24h := trading.TradingLast24h{trdconst.LAST24H, last24h}
		trdLast24hJson, err := json.Marshal(trdLast24h)
		if err!=nil{
			log.Error("error when marshal last24h", err)
		}

		log.Println(trdLast24h)
		cache.SetCacheByTopicAndType(topic, trdconst.LAST24H, trdLast24h)

		websocket.BroadcastMessage(topic, string(trdLast24hJson))

	}


}


func ListenTradingHistory(publisher string, topic string, zmqKey string, clients *daos.Clients){

	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
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

		log.Println(trdListHistory)
		cache.SetCacheByTopicAndType(topic, trdconst.TRADINGHISTORY, trdListHistory)

		websocket.BroadcastMessage(topic, string(trdListHistoryJson))

	}
}

func ListenOpenOrder(publisher string, topic string, zmqKey string, clients *daos.Clients){

	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
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
		log.Println(trdListOpenOrder)
		websocket.SendMessageToUser(topic, username, string(trdListOpenOrderJson))


	}
}

func ListenOrderHistory(publisher string, topic string, zmqKey string, clients *daos.Clients){
	const SIZE = 5

	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)
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
		log.Println(trdListOrderHistory)
		websocket.SendMessageToUser(topic, username, string(trdListOrderHistoryJson))

	}
}