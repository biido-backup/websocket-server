package main

import (
	"fmt"
	"github.com/go-zeromq/zmq4"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"time"
	"context"
)

var clients map[string] map[string] sockjs.Session

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Rate struct{
	MainCurrency		string 		`json:"mainCurrency"`
	PivotCurrency		string 		`json:"pivotCurrency"`
}

func main(){

	clients = make(map[string] map[string] sockjs.Session)
	//go clientHub()

	//go socketListener("BTC-IDR")
	//go socketListener("XRP-IDR")
	go socketListener("BTC-IDR:ORDER_BOOK")

	handler := sockjs.NewHandler("/echo", sockjs.DefaultOptions, echoHandler)

	log.Fatal(http.ListenAndServe(":18000", handler))


	fmt.Println("Done Main")
}

func echoHandler(session sockjs.Session) {
	str := string("connection from server (" + session.ID() + ") : OPEN")

	//register client
	subscribeClientToTopic("BTC-IDR", session)
	session.Send("start subscribe to : BTC-IDR")
	fmt.Println(clients)
	fmt.Println(str)
	for {
		if msg, err := session.Recv(); err == nil {
			fmt.Println(msg)

			unsubscribeClientToAllTopic(session.ID())
			subscribeClientToTopic(msg, session)

			str := string("subscribe to : "+msg)
			session.Send(str)

			continue
		}
		str := string("connection from server (" + session.ID() + ") : CLOSED")
		fmt.Println(str)

		//remove client
		unsubscribeClientToAllTopic(session.ID())
		fmt.Println(clients)
		break
	}
}


func socketListener(rate string){
	clients[rate] = make(map[string] sockjs.Session)

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
		//var message Rate
		//_ = json.Unmarshal(msg.Frames[1], &message)
		//broadcastMessage(rate, "mainCurrency : "+message.MainCurrency)
		fmt.Println(msg)
	}

}

func unsubscribeClientToAllTopic(sessionID string){
	for topic, _ := range(clients) {
		delete(clients[topic], sessionID)
	}
}

func subscribeClientToTopic(topic string, session sockjs.Session){
	fmt.Println("subscribe")
	clients[topic][session.ID()] = session

}


func broadcastMessage(topic string, str string){


	start := time.Now()
	for _, session := range(clients[topic]) {
		session.Send(str)
	}
	fmt.Println("Time : "+time.Since(start).String())
	log.Println(clients)


}
