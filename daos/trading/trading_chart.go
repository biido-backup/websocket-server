package trading

import "websocket-server/daos"

type TradingChart struct {
	Type		string					`json:"type"`
	Payload		[]daos.TradingChart		`json:"payload"`
}
