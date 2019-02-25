package main

import (
	"log"
	"sync"
	"time"
	"websocket-server/util/config"
	"websocket-server/util/database"
)

var mu *sync.RWMutex
var i = make(map[string] int)

func main() {
	config.LoadConfig();
	database.ConnectDbPostgres()

	//var listOrderHistory trading.ListOrderHistory
	//service.GetListOrderHistoryByMemberIdAndRateAndOffsetAndLimit(&listOrderHistory, 99, "BTC/IDR", 0, 10)
	//
	//log.Println(listOrderHistory.Size)
	//log.Println(listOrderHistory.OrderHistories)

	//var listOpenOrder trading.ListOpenOrder
	//service.GetListOpenOrderByMemberIdAndRate(&listOpenOrder, 99, "BTC/IDR")
	//
	//log.Println(listOpenOrder)

	//var listHistory	trading.ListHistory
	//service.GetLastTradingHistoryBySize(&listHistory, "LTC/IDR", 10)
	//
	//log.Println(listHistory)

	//var l trading.Last24h
	//service.GetLast24HTransactionByRate(&l, "LTC/IDR")
	//
	//log.Println(l)

	//log.Println("TESS")
	//var openOrders trading.OpenOrders
	//err := service.GetOpenOrdersByUsernameAndRate(&openOrders, "agung.wb", "BTC/IDR")
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(openOrders)

	//log.Println("TESS")
	//var orderHistories trading.OrderHistories
	//err := service.GetOrderHistoriesByUsernameAndRateAndOffsetAndLimit(&orderHistories, "agung.wb", "BTC/IDR", 0, 10)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(orderHistories.Payload)
	//log.Println(orderHistories.Type)
	//log.Println(orderHistories.Size)

	mu = &sync.RWMutex{}

	i["tes"] = 0
	from := 1000000
	for ii:=0; ii<from; ii++ {
		go inc()
	}

	log.Println(i["tes"])
	time.Sleep(time.Minute)
}


func inc(){
	mu.Lock()
	i["tes"] += 1
	log.Println(i["tes"])
	mu.Unlock()
}