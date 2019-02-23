package trading

import (
	"websocket-server/const/trd"
	"websocket-server/daos"
)

type TradingChart struct {
	Type		string				`json:"type"`
	Payload		[]daos.Chart		`json:"payload"`
}

func CreateTradingChart(subscriber daos.Subscriber) TradingChart {
	unitOfTime := ""
	if subscriber.Interval == "1m" {
		unitOfTime = "MINUTE"
	} else if subscriber.Interval == "1h" {
		unitOfTime = "HOUR"
	} else if subscriber.Interval == "1d" || subscriber.Interval == "1w" {
		unitOfTime = "DAY"
	} else if subscriber.Interval == "1M"{
		unitOfTime = "MONTH"
	}

	return TradingChart{trdconst.TRADINGCHART, daos.GetChartList(subscriber.Topic, unitOfTime)}
}