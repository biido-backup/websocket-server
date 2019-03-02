package trading

import (
	"websocket-server/const/trd"
	"websocket-server/daos"
)

type TradingChart struct {
	Type		string				`json:"type"`
	Payload		[]daos.Chart		`json:"payload"`
}

func CreateTradingChart(subscriber daos.WebsocketRequest) TradingChart {
	unitOfTime := ""
	if subscriber.Interval == "1m" {
		unitOfTime = "1M"
	} else if subscriber.Interval == "5m" {
		unitOfTime = "5M"
	} else if subscriber.Interval == "15m" {
		unitOfTime = "15M"
	} else if subscriber.Interval == "30m" {
		unitOfTime = "30M"
	} else if subscriber.Interval == "1h" {
		unitOfTime = "1H"
	} else if subscriber.Interval == "6h" {
		unitOfTime = "6H"
	} else if subscriber.Interval == "12h" {
		unitOfTime = "12H"
	} else if subscriber.Interval == "1d" {
		unitOfTime = "1D"
	} else if subscriber.Interval == "1w" {
		unitOfTime = "1W"
	}  else if subscriber.Interval == "1M"{
		unitOfTime = "1MO"
	}

	return TradingChart{trdconst.TRADINGCHART, daos.GetChartList(subscriber.Topic, unitOfTime)}
}