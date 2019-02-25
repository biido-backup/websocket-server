package main

import (
	"log"
	"websocket-server/daos/trading"
	"websocket-server/service"
	"websocket-server/util/config"
	"websocket-server/util/database"
)

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

	var l trading.Last24h
	service.GetLast24HTransactionByRate(&l, "LTC/IDR")

	log.Println(l)

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

}
