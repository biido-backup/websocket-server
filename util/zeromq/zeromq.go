package zeromq

import (
	"encoding/json"
	"github.com/go-zeromq/zmq4"
	"log"
	"context"
	"websocket-server/daos"
	"websocket-server/util/websocket"
)


func Listen(topic string, clients *daos.Clients){
	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	//dial
	err := sub.Dial("tcp://localhost:5563")
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, topic)
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Fatalf("could not receive message: %v", err)
		}

		//log.Println(msg)
		//
		//orderbook := daos.OrderBookFromJSONZeroMQ(msg.Bytes())

		var orderbook daos.OrderBookZeroMQ;
		json.Unmarshal(msg.Bytes(), &orderbook);

		log.Println(orderbook)

		//msg.Frames[0] --> topic
		//msg.Frames[1] --> message
		var message daos.Rate
		_ = json.Unmarshal(msg.Frames[1], &message)
		websocket.BroadcastMessage(topic, "mainCurrency : "+message.MainCurrency)

	}

}


