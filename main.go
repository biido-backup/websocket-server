package main

import (
	"github.com/sirupsen/logrus"
	"sync"
	"websocket-server/daos"
	"websocket-server/service"
	"websocket-server/util/config"
	"websocket-server/util/database"
	"websocket-server/util/kafka"
	"websocket-server/util/websocket"
	"websocket-server/util/zeromq"
)

var log = logrus.New()
var clients *daos.Clients
var mutex *sync.Mutex

func main(){
	config.LoadConfig()
	database.ConnectDbPostgres()

	clients = daos.CreateClients()
	log.Println(*clients)

	mutex = &sync.Mutex{}
	daos.CreateTradingCharts()

	walletTypeList := service.GetAllWalletType()
	rateList := make([]string, 0, 1)
	for _, pivot := range walletTypeList{
		if pivot.FlagActive == 1 && pivot.PivotPriority.Valid {
			for _, main := range walletTypeList{
				if main.FlagActive == 1 && (!main.PivotPriority.Valid || main.PivotPriority.Int64 > pivot.PivotPriority.Int64){
					rateList = append(rateList,  main.Code+ "-" +pivot.Code)
				}
			}
		}
	}
	rateList = []string{"BTC-IDR"}
	unitOfTimeList := []string{"MINUTE"}
	for _, rate := range rateList {
		for _, unitOfTime := range unitOfTimeList{
			initMapList := make(map[string][]string)
			initMapOffset := make(map[string]int64)
			initMapList[unitOfTime], initMapOffset[unitOfTime] = kafka.GetCandleStick(rate, unitOfTime)

			daos.Charts.ListMap[rate] = initMapList
			daos.Charts.OffsetMap[rate] = initMapOffset

			go kafka.MaintenanceTradingChartArray(rate, unitOfTime, daos.Charts, mutex)
		}
	}

	go zeromq.Listen("BTC-IDR:ORDER_BOOK", clients)
	//go zeromq.Listen("XRP-IDR", clients)

	err := websocket.ServeSocket(clients)
	if err != nil {
		log.Error(err)
		log.Fatal(err)
	}

}




