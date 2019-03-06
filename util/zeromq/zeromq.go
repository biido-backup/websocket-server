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

var HasConnectedZeroMQMatchingEngine = false
var HasConnectedZeroMQTradingBroker = false

func Listen(tradingRateList []daos.Rate)  {
	done := make(chan bool)

	cronZMQMatchingEngine := cron.New()
	cronZMQMatchingEngine.AddFunc("0-59/5 * * * * *", func() {
		if !HasConnectedZeroMQMatchingEngine {
			log.Debug("Reconnect to zeromq matching engine client")
			for _, tradingRate := range(tradingRateList){
				ListenMatchingEngine(tradingRate.StringDash())
			}
		}
	})
	cronZMQMatchingEngine.Start()

	cronZMQTradingBroker := cron.New()
	cronZMQTradingBroker.AddFunc("0-59/5 * * * * *", func() {
		if !HasConnectedZeroMQTradingBroker {
			log.Debug("Reconnect to zeromq trading broker client")
			for _, tradingRate := range(tradingRateList){
				ListenTradingBroker(tradingRate.StringDash())
			}
		}
	})
	cronZMQTradingBroker.Start()

	<- done
	
}

func ListenMatchingEngine(topic string){
	var matchingEngineAddr = viper.GetString("zeromq.publisher.matching-engine")
	var suffixOrderBook = viper.GetString("zeromq.key.suffix.orderbook")

	go ListenOrderBook(matchingEngineAddr, topic, topic+suffixOrderBook, &daos.MyClients)
}

func ListenTradingBroker(topic string){
	var tradingBrokerAddr = viper.GetString("zeromq.publisher.trading-broker")

	var suffixLast24h = viper.GetString("zeromq.key.suffix.last24h")
	var suffixTradingHistory = viper.GetString("zeromq.key.suffix.tradinghistory")
	var suffixOpenOrder = viper.GetString("zeromq.key.suffix.openorder")
	var suffixOrderHistory = viper.GetString("zeromq.key.suffix.orderhistory")

	go ListenLast24h(tradingBrokerAddr, topic, topic+suffixLast24h, &daos.MyClients)
	go ListenTradingHistory(tradingBrokerAddr, topic, topic+suffixTradingHistory, &daos.MyClients)
	go ListenOpenOrder(tradingBrokerAddr, topic, topic+suffixOpenOrder, &daos.MyClients)
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

		log.Debug(orderbook)

	}
}

func ListenOrderBook(publisher string, topic string, zmqKey string, clients *daos.Clients){
	clients.SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQMatchingEngine = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQMatchingEngine = false
		return
	}

	HasConnectedZeroMQMatchingEngine = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQMatchingEngine = false
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

		log.Debug(orderbook)
		cache.SetCacheByTopicAndType(topic, trdconst.ORDERBOOK, trdOrderbook)
		websocket.BroadcastMessageWithTopic(topic, string(trdOrderbookJson))
	}
}


func ListenLast24h(publisher string, topic string, zmqKey string, clients *daos.Clients){
	clients.SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	HasConnectedZeroMQTradingBroker = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQTradingBroker = false
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


func ListenTradingHistory(publisher string, topic string, zmqKey string, clients *daos.Clients){
	clients.SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	HasConnectedZeroMQTradingBroker = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQTradingBroker = false
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

func ListenOpenOrder(publisher string, topic string, zmqKey string, clients *daos.Clients){
	clients.SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	HasConnectedZeroMQTradingBroker = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQTradingBroker = false
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

func ListenOrderHistory(publisher string, topic string, zmqKey string, clients *daos.Clients){
	const SIZE = 5

	clients.SetTopic(topic)
	sub := zmq4.NewSub(context.Background())

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Error("could not dial publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Error("could not subscribe publisher "+publisher, err)

		HasConnectedZeroMQTradingBroker = false
		return
	}

	HasConnectedZeroMQTradingBroker = true
	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Error("could not receive message", err)

			sub.Close()
			HasConnectedZeroMQTradingBroker = false
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