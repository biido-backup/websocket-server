package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
	"websocket-server/cache"
	"websocket-server/daos"
	"websocket-server/engine"
	"websocket-server/util/config"
	"websocket-server/util/database"
	"websocket-server/util/logger"
	"websocket-server/util/redis"
	"websocket-server/util/websocket"
	"websocket-server/util/zeromq"
)

var log *logger.CustomLog

var clients *daos.Clients
var mutex *sync.Mutex

func init(){
	config.LoadConfig()
	log = logger.CreateLog("main")
	fmt.Println("Logging Level : "+log.Level)
	redis.ConnectRedis()
	database.ConnectDbPostgres()
	cache.InitCache()
}

func main(){

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
		cache.FillCacheByRate(tradingRate)
		zeromq.Listen(tradingRate.StringDash(), &clients)
		time.Sleep(time.Millisecond)
	}
	engine.ProcessTradingChart(tradingRateList)

	err := websocket.ServeSocket(&clients)
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}

}




