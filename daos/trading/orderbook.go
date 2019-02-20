package trading

import "websocket-server/daos"

type Orderbook struct {
	Type 		string 						`json:"type"`
	Payload		daos.OrderBookZeroMQ		`json:"payload"`
}


