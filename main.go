package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"websocket-server/daos"
	"websocket-server/util/config"
	"websocket-server/util/redis"
	"websocket-server/util/websocket"
	"websocket-server/util/zeromq"
)

var log = logrus.New()
var clients *daos.Clients
var mutex *sync.Mutex


func main(){
	config.LoadConfig()
	redis.ConnectRedis()

	var clients daos.Clients
	clients = daos.CreateClients()

	log.Println(clients)

	rateList := viper.GetString("redis.trading.key")

	tradingRateJson := redis.GetValueByKey(rateList)
	tradingRateList := make([]daos.Rate, 0, 1)
	json.Unmarshal(tradingRateJson, &tradingRateList)

	for _, tradingRate := range(tradingRateList){
		log.Println(tradingRate)
		clients.SetTopic(tradingRate.StringDash())
		go zeromq.Listen(tradingRate.StringDash()+":"+"ORDER_BOOK", &clients)
		//time.Sleep(time.Millisecond)
	}

	err := websocket.ServeSocket(&clients)
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}

}




