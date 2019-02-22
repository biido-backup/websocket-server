package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
)

var log = logrus.New()

var clients *daos.Clients

func ServeSocket(cls *daos.Clients) error{

	clients = cls

	port := viper.GetString("sockjs.port")
	path := viper.GetString("sockjs.path")

	log.Println(port)
	log.Println(path)

	handler := sockjs.NewHandler(path, sockjs.DefaultOptions, SockjsHandler)
	return http.ListenAndServe(":"+port, handler)
}

func SockjsHandler(session sockjs.Session) {
	str := string("connection from server (" + session.ID() + ") : OPEN")

	//register client
	//subscribeClientToTopic("BTC-IDR", session)
	//session.Send("start subscribe to : BTC-IDR")
	fmt.Println(clients)
	fmt.Println(str)

	for {
		if msg, err := session.Recv(); err == nil {

			//log.Println(msg)

			var subscriber daos.Subscriber
			err := json.Unmarshal([]byte(msg), &subscriber)

			//log.Println(subscriber)

			if err != nil {
				log.Error(err)
				continue
			}

			unsubscribeClientToAllTopic(session.ID())
			subscribeClientToTopic(subscriber, session)

			str := string("subscribe to : "+subscriber.Topic)
			session.Send(str)

			tradingChart := trading.TradingChart{trdconst.TRADINGCHART, daos.GetTradingChartListMap()[subscriber.Topic]["MINUTE"]}
			tradingChartJson, _ := json.Marshal(tradingChart)
			session.Send(string(tradingChartJson))

			continue
		}
		str := string("connection from server (" + session.ID() + ") : CLOSED")
		log.Println(str)

		//remove client
		unsubscribeClientToAllTopic(session.ID())
		log.Println(clients.Clients)
		log.Println(clients.ClientSessions)
		break
	}
}

func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(clients.ClientSessions) {
		clients.RemoveSubscriber(topic, sessionID)
	}
}

func subscribeClientToTopic(subscriber daos.Subscriber, session sockjs.Session){
	log.Println("subscribe")

	clients.AddSubscriber(subscriber, session)

	log.Println(clients.Clients)
	log.Println(clients.ClientSessions)
	log.Println(clients.Intervals)
	log.Println(clients.IntervalSessions)
}

func BroadcastMessage(topic string, str string){
	//start := time.Now()
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
