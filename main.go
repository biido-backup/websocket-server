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

//var clients *daos.Clients
var mutex *sync.Mutex

func init(){
	_, _ = time.LoadLocation("Asia/Jakarta")
	config.LoadConfig()
	log = logger.CreateLog("main")
	fmt.Println("Logging Level : "+log.Level)
	redis.ConnectRedis()
	database.ConnectDbPostgres()
	cache.InitCache()
}

func main(){

	//var clients daos.Clients
	//clients = daos.CreateClients()

	daos.CreateClients()

	log.Println(daos.MyClients)

	rateList := viper.GetString("redis.trading.key")

	tradingRateJson := redis.GetValueByKey(rateList)
	tradingRateList := make([]daos.Rate, 0, 1)
	json.Unmarshal(tradingRateJson, &tradingRateList)

	for _, tradingRate := range(tradingRateList){
		daos.MyClients.SetTopic(tradingRate.StringDash())
		cache.FillCacheByRate(tradingRate)
	}

	go zeromq.Listen()
	go engine.ProcessTradingChart(tradingRateList)

	err := websocket.ServeSocket()
	if err != nil {
		log.Fatal(err)
	}


}




