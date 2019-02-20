package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"github.com/sirupsen/logrus"
	"net/http"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"context"
	"websocket-server/daos"
)

var log = logrus.New()

var clients map[string] map[string] map[string] sockjs.Session  //topic -> username -> sessionID
var sessions map[string] map[string] string //topic -> sessionId -> username

func main(){

	clients = make(map[string] map[string] map[string] sockjs.Session)
	sessions = make(map[string] map[string] string)
	//go clientHub()

	go socketListener("BTC-IDR")
	go socketListener("XRP-IDR")

	handler := sockjs.NewHandler("/echo", sockjs.DefaultOptions, sockjsHandler)

	log.Fatal(http.ListenAndServe(":18000", handler))


	fmt.Println("Done Main")
}

func sockjsHandler(session sockjs.Session) {
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
		log.Println(clients)
		log.Println(sessions)
		break
	}
}


func socketListener(rate string){
	clients[rate] = make(map[string] map[string] sockjs.Session)
	sessions[rate] = make(map[string] string)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial("tcp://localhost:5563")
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, rate)
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Fatalf("could not receive message: %v", err)
		}
		//msg.Frames[0] --> topic
		//msg.Frames[1] --> message
		var message daos.Rate
		_ = json.Unmarshal(msg.Frames[1], &message)
		broadcastMessage(rate, "mainCurrency : "+message.MainCurrency)

	}

}

func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(sessions) {
		username := sessions[topic][sessionID]

		delete(clients[topic][username], sessionID)
		if len(clients[topic][username]) == 0 {
			delete(clients[topic], username)
		}
		delete(sessions[topic], sessionID)
	}
}

//func unsubscribeClientToAllTopic(username string, sessionID string){
//	log.Println("unsubscribe")
//	for topic, _ := range(clients) {
//
//		log.Println(topic)
//		log.Println(username)
//		log.Println(sessionID)
//		log.Println(clients[topic][username][sessionID])
//
//		if (clients[topic][username] != nil){
//			log.Println(topic)
//			log.Println(username)
//			log.Println(sessionID)
//
//			delete(clients[topic][username], sessionID)
//			if len(clients[topic][username]) == 0 {
//				delete(clients[topic], username)
//			}
//			delete(sessions[topic], sessionID)
//		}
//	}
//}

func subscribeClientToTopic(topic string, username string, session sockjs.Session){
	log.Println("subscribe")
	if clients[topic][username] == nil{
		clients[topic][username] = make(map[string] sockjs.Session)
	}

	clients[topic][username][session.ID()] = session
	sessions[topic][session.ID()] = username

	log.Println(clients)
	log.Println(sessions)
}


func broadcastMessage(topic string, str string){

	//start := time.Now()
	for _, username := range(clients[topic]) {
		for _, session := range(username){
			session.Send(str)
		}

	}
	//log.Println("Time : "+time.Since(start).String())
	//log.Println(clients)
	//log.Println(sessions)

}

func sendMessageToUser(topic string, username string, str string){

	for _, session := range(clients[topic][username]) {
			session.Send(str)
	}

}
