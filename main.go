package main

import (
	"github.com/sirupsen/logrus"
	"websocket-server/daos"
	"websocket-server/util/config"
	"websocket-server/util/database"
)

var log = logrus.New()

var clients *daos.Clients


func main(){
	config.LoadConfig();

	//clients = daos.CreateClients()
	//log.Println(*clients)
	database.ConnectDbPostgres()


	//walletTypeList, err := service.GetAllWalletType()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, walletType := range(walletTypeList){
	//	log.Println(walletType)
	//}

	//go zeromq.Listen("BTC-IDR:ORDER_BOOK", clients)
	////go zeromq.Listen("XRP-IDR", clients)
	//
	//err := websocket.ServeSocket(clients)
	//if err != nil {
	//	log.Error(err)
	//	log.Fatal(err)
	//}

}




