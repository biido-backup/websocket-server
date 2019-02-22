package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"websocket-server/daos"
	"websocket-server/util/config"
	"websocket-server/util/redis"
	"websocket-server/util/websocket"
)

var log = logrus.New()

var clients *daos.Clients


func main(){
	config.LoadConfig()
	redis.ConnectRedis()

	var clients daos.Clients
	clients = daos.CreateClients()

	log.Println(clients)

	rate := viper.GetString("redis.trading.key")

	tradingRateJson := redis.GetValueByKey(rate)
	tradingRateList := make([]daos.Rate, 0, 1)
	json.Unmarshal(tradingRateJson, &tradingRateList)

	for _, tradingRate := range(tradingRateList){
		log.Println(tradingRate)
		clients.SetTopic(tradingRate.StringDash())
	}

	err := websocket.ServeSocket(&clients)
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}

}




