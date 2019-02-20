package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"github.com/gorilla/websocket"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
	"strconv"
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


func main2(){
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":18000", nil)
}

func main(){

	clients = make(map[string] map[string] sockjs.Session)
	//go clientHub()

	go socketListener("BTC-IDR")
	go socketListener("XRP-IDR")

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

func echoHandler2(session sockjs.Session) {
	i := 0
	for {
		//if _, err := session.Recv(); err == nil {
		str := string("from server (" + session.ID() + ") : " + strconv.Itoa(i))
		session.Send(str)
		fmt.Println(str)
		i += 1
		time.Sleep(time.Second)


		//	continue
		//}
		//fmt.Println("oi")
	}

	session.Recv()
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
		var message Rate
		_ = json.Unmarshal(msg.Frames[1], &message)
		broadcastMessage(rate, "mainCurrency : "+message.MainCurrency)

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

func clientHub(){
	counter := 0
	for {
		str := "Broadcast from server : "+strconv.Itoa(counter)
		fmt.Println(str)
		broadcastMessage("BTC-IDR", str)
		time.Sleep(time.Second * 3)
		counter += 1
	}
}

func broadcastMessage(topic string, str string){


	start := time.Now()
	for _, session := range(clients[topic]) {
		session.Send(str)
	}
	fmt.Println("Time : "+time.Since(start).String())
	log.Println(clients)


}

