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

//var clients = &daos.MyClients

//func ServeSocket(cls *daos.Clients) error{
func ServeSocket() error{

	//clients = cls
	//clients = &daos.MyClients

	port := viper.GetString("sockjs.port")
	path := viper.GetString("sockjs.path")


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

			var request daos.WebsocketRequest
			err := json.Unmarshal([]byte(msg), &request)

			if err != nil {
				log.Error("Failed to deserialize message", err)
				continue
			}

			rate := daos.GetRateFromStringDash(request.Topic)

			log.Debug(request)

			if request.Method == "subscribe"{
				unsubscribeClientToAllTopic(session.ID())
				subscribeClientToTopic(request, session)

				str := string("subscribe to : "+ request.Topic)
				session.Send(str)

				//Candle Stick
				tradingChart := trading.CreateTradingChart(request)
				tradingChartJson, _ := json.Marshal(tradingChart)
				session.Send(string(tradingChartJson))

				//Orderbook
				orderBook := cache.GetCacheByTopicAndType(request.Topic, trdconst.ORDERBOOK).(trading.Orderbook)
				orderBookJson, _ := json.Marshal(orderBook)
				session.Send(string(orderBookJson))

				//TradingHistory
				tradingHistory := cache.GetCacheByTopicAndType(request.Topic, trdconst.TRADINGHISTORY).(trading.TradingListHistory)
				tradingHistoryJson, _ := json.Marshal(tradingHistory)
				session.Send(string(tradingHistoryJson))

				//Last24H
				//last24h := cache.GetCacheByTopicAndType(request.Topic, trdconst.LAST24H).(trading.TradingLast24h)
				//last24hJson, _ := json.Marshal(last24h)
				//session.Send(string(last24hJson))

				//AllRatesLast24H
				allRatesLast24h := cache.GetCacheByType(trdconst.LAST24H).(map[string]trading.TradingLast24h)
				for _, last24h := range(allRatesLast24h){
					log.Debug(last24h)
					last24hJson, _ := json.Marshal(last24h)
					session.Send(string(last24hJson))
				}

				//OrderHistory
				const offset uint64 = 0
				const limit int = 5
				var orderHistories trading.OrderHistories
				err = service.GetOrderHistoriesByUsernameAndRateAndOffsetAndLimit(&orderHistories, request.Username, rate.StringSlah(), offset, limit)
				if err != nil {
					log.Println(err)
				}
				orderHistories.Type = trdconst.ORDERHISTORY
				orderHistoriesJson, _ := json.Marshal(orderHistories)
				session.Send(string(orderHistoriesJson))

				//OpenOrder
				var openOrders trading.OpenOrders
				err = service.GetOpenOrdersByUsernameAndRate(&openOrders, request.Username, rate.StringSlah())
				if err != nil {
					log.Error(err)
				}
				openOrders.Type = trdconst.OPENORDER
				openOrdersJson, _ := json.Marshal(openOrders)
				session.Send(string(openOrdersJson))

				continue

			} else if request.Method == "reload_openorder" {
				//OpenOrder
				log.Debug("Reload Open Order")
				if checkIfSubscribed(request.Topic, request.Username, session.ID()){
					log.Debug("Reload Open Order: MASUK")
					var openOrders trading.OpenOrders
					err = service.GetOpenOrdersByUsernameAndRate(&openOrders, request.Username, rate.StringSlah())
					if err != nil {
						log.Error(err)
					}
					openOrders.Type = trdconst.OPENORDER
					openOrdersJson, _ := json.Marshal(openOrders)
					SendMessageToUser(request.Topic, request.Username, string(openOrdersJson))
				}
				continue
			}
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

func checkIfSubscribed(topic string, username string, sessionId string) bool{
	return daos.MyClients.CheckIfSessionExistsByTopicAndUsername(topic, username, sessionId)
}

func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(daos.MyClients.ClientSessions) {
		daos.MyClients.RemoveSubscriber(topic, sessionID)
	}
}

func subscribeClientToTopic(request daos.WebsocketRequest, session sockjs.Session){
	daos.MyClients.AddSubscriber(request, session)
	//log.Println(clients.Clients)
	//log.Println(clients.ClientSessions)
	//log.Println(clients.Intervals)
	//log.Println(clients.IntervalSessions)
}

func BroadcastMessageToAll(str string){
	//start := time.Now()
	var c map[string] map[string] map[string] sockjs.Session
	c = daos.MyClients.GetAllClients()
	for _,topic := range(c){
		for _, username := range(topic) {
			for _, session := range(username){
				session.Send(str)
			}
		}
	}
}

func BroadcastMessageWithTopic(topic string, str string){
	//start := time.Now()
	var c map[string] map[string] sockjs.Session
	c = daos.MyClients.GetAllClientsByTopic(topic)
	for _, username := range(c) {
		for _, session := range(username){
			session.Send(str)
		}
	}
}

func BroadcastMessageWithInterval(topic string, interval string, str string){
	var sessions map[string] sockjs.Session
	sessions = daos.MyClients.GetListSessionByTopicAndInterval(topic, interval)
	for _, session := range(sessions) {
		session.Send(str)
	}
}

func SendMessageToUser(topic string, username string, str string){
	var sessions map[string] sockjs.Session
	sessions = daos.MyClients.GetListSessionByTopicAndUsername(topic, username)
	for _, session := range(sessions) {
		session.Send(str)
	}
}