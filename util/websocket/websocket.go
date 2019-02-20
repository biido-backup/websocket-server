package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"websocket-server/daos"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"github.com/spf13/viper"
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

			log.Println(msg)

			var subscriber daos.Subscriber
			err := json.Unmarshal([]byte(msg), &subscriber)

			log.Println(subscriber)

			if err != nil {
				log.Error(err)
				continue
			}

			unsubscribeClientToAllTopic(session.ID())
			subscribeClientToTopic(subscriber.Topic, subscriber.Username, session)



			str := string("subscribe to : "+subscriber.Topic)
			session.Send(str)

			continue
		}
		str := string("connection from server (" + session.ID() + ") : CLOSED")
		log.Println(str)

		//remove client
		unsubscribeClientToAllTopic(session.ID())
		log.Println(clients.Clients)
		log.Println(clients.Sessions)
		break
	}
}


func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(clients.Sessions) {
		clients.DeleteClient(topic, sessionID)
	}
}


func subscribeClientToTopic(topic string, username string, session sockjs.Session){
	log.Println("subscribe")

	clients.AddClient(topic, username , session)

	log.Println(clients.Clients)
	log.Println(clients.Sessions)
}


func BroadcastMessage(topic string, str string){


	//start := time.Now()
	for _, username := range(clients.Clients[topic]) {
		for _, session := range(username){
			session.Send(str)
		}

	}
	//log.Println("Time : "+time.Since(start).String())
	//log.Println(clients)
	//log.Println(sessions)

}

func SendMessageToUser(topic string, username string, str string){

	for _, session := range(clients.Clients[topic][username]) {
		session.Send(str)
	}

}
