package zeromq

import (
	"encoding/json"
	"github.com/go-zeromq/zmq4"
	"github.com/sirupsen/logrus"
	"context"
	"github.com/spf13/viper"
	"strings"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/websocket"
)

var log = logrus.New()



func Listen(zmqKey string, clients *daos.Clients){

	topic := strings.Split(zmqKey,":")[0]
	clients.SetTopic(topic)

	sub := zmq4.NewSub(context.Background())
	defer sub.Close()

	var publisher = viper.GetString("zeromq.publisher")

	//dial
	err := sub.Dial(publisher)
	if err != nil {
		log.Fatalf("could not dial: %v", err)
	}

	//subscribe
	err = sub.SetOption(zmq4.OptionSubscribe, zmqKey)
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	for {
		// Read envelope
		msg, err := sub.Recv()
		if err != nil {
			log.Fatalf("could not receive message: %v", err)
		}
		//
		orderbook := daos.OrderBookFromJSONZeroMQ(msg.Frames[1])

		log.Println(orderbook)

		trdOrderbook := trading.Orderbook{trdconst.ORDERHISTORY, orderbook}
		trdOrderbookJson, err := json.Marshal(trdOrderbook)
		if err!=nil{
			log.Error(err)
		}


		//msg.Frames[0] --> zmqKey
		//msg.Frames[1] --> message
		var message daos.Rate
		_ = json.Unmarshal(msg.Frames[1], &message)
		websocket.BroadcastMessage(topic, string(trdOrderbookJson))

	}

}


