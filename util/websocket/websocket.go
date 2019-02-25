package websocket

import (
	"encoding/json"
	"github.com/spf13/viper"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
	"websocket-server/cache"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/service"
	"websocket-server/util/logger"
)

var log = logger.CreateLog("websocket")

var clients *daos.Clients

func ServeSocket(cls *daos.Clients) error{

	clients = cls

	port := viper.GetString("sockjs.port")
	path := viper.GetString("sockjs.path")

	log.Debug(port)
	log.Debug(path)

	handler := sockjs.NewHandler(path, sockjs.DefaultOptions, SockjsHandler)
	return http.ListenAndServe(":"+port, handler)
}

func SockjsHandler(session sockjs.Session) {
	str := string("connection from server (" + session.ID() + ") : OPEN")

	//register client
	//subscribeClientToTopic("BTC-IDR", session)
	//session.Send("start subscribe to : BTC-IDR")
	//fmt.Println(clients)
	log.Debug(str)

	for {
		if msg, err := session.Recv(); err == nil {

			log.Debug(msg)

			var subscriber daos.Subscriber
			err := json.Unmarshal([]byte(msg), &subscriber)

			if err != nil {
				log.Error("Failed to deserialize message", err)
				continue
			}

			rate := daos.GetRateFromStringDash(subscriber.Topic)

			log.Debug(subscriber)

			unsubscribeClientToAllTopic(session.ID())
			subscribeClientToTopic(subscriber, session)

			str := string("subscribe to : "+subscriber.Topic)
			session.Send(str)

			//Candle Stick
			tradingChart := trading.CreateTradingChart(subscriber)
			tradingChartJson, _ := json.Marshal(tradingChart)
			session.Send(string(tradingChartJson))

			//Orderbook
			orderBook := cache.GetCacheByTopicAndType(subscriber.Topic, trdconst.ORDERBOOK).(trading.Orderbook)
			orderBookJson, _ := json.Marshal(orderBook)
			session.Send(string(orderBookJson))

			//TradingHistory
			tradingHistory := cache.GetCacheByTopicAndType(subscriber.Topic, trdconst.TRADINGHISTORY).(trading.TradingListHistory)
			tradingHistoryJson, _ := json.Marshal(tradingHistory)
			session.Send(string(tradingHistoryJson))

			//Last24H
			last24h := cache.GetCacheByTopicAndType(subscriber.Topic, trdconst.LAST24H).(trading.TradingLast24h)
			last24hJson, _ := json.Marshal(last24h)
			session.Send(string(last24hJson))

			//OrderHistory
			var orderHistories trading.OrderHistories
			err = service.GetOrderHistoriesByUsernameAndRateAndOffsetAndLimit(&orderHistories, subscriber.Username, rate.StringSlah(), 0, 10)
			if err != nil {
				log.Println(err)
			}
			orderHistories.Type = trdconst.ORDERHISTORY
			orderHistoriesJson, _ := json.Marshal(orderHistories)
			session.Send(string(orderHistoriesJson))

			//OpenOrder
			var openOrders trading.OpenOrders
			err = service.GetOpenOrdersByUsernameAndRate(&openOrders, subscriber.Username, rate.StringSlah())
			if err != nil {
				log.Println(err)
			}

			continue
		}
		str := string("connection from server (" + session.ID() + ") : CLOSED")
		log.Debug(str)

		//remove client
		unsubscribeClientToAllTopic(session.ID())
		//log.Println(clients.Clients)
		//log.Println(clients.ClientSessions)

		break
	}
}

func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(clients.ClientSessions) {
		clients.RemoveSubscriber(topic, sessionID)
	}
}

func subscribeClientToTopic(subscriber daos.Subscriber, session sockjs.Session){
	log.Debug("subscribe : ")
	log.Debug(subscriber)

	clients.AddSubscriber(subscriber, session)

	log.Println(clients.Clients)
	//log.Println(clients.ClientSessions)
	//log.Println(clients.Intervals)
	//log.Println(clients.IntervalSessions)
}

func BroadcastMessage(topic string, str string){
	//start := time.Now()
	log.Debug("Broadcast message with topic : "+topic)
	var c map[string] map[string] sockjs.Session
	c = clients.GetAllClientsByTopic(topic)
	for _, username := range(c) {
		for _, session := range(username){
			session.Send(str)
		}

	}
	//log.Println("Time : "+time.Since(start).String())
	//log.Println(clients)
	//log.Println(sessions)
}

func BroadcastMessageWithInterval(topic string, interval string, str string){
	var sessions map[string] sockjs.Session
	sessions = clients.GetListSessionByTopicAndInterval(topic, interval)
	for _, session := range(sessions) {
		session.Send(str)
	}
}

func SendMessageToUser(topic string, username string, str string){

	var sessions map[string] sockjs.Session
	sessions = clients.GetListSessionByTopicAndUsername(topic, username)
	for _, session := range(sessions) {
		session.Send(str)
	}

}
