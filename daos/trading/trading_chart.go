package trading

import (
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/reqpayload"
)

type TradingChart struct {
	Type		string				`json:"type"`
	Payload		[]daos.Chart		`json:"payload"`
}

func CreateTradingChart(payload reqpayload.Trading) TradingChart {
	unitOfTime := ""
	if payload.Interval == "1m" {
		unitOfTime = "1M"
	} else if payload.Interval == "5m" {
		unitOfTime = "5M"
	} else if payload.Interval == "15m" {
		unitOfTime = "15M"
	} else if payload.Interval == "30m" {
		unitOfTime = "30M"
	} else if payload.Interval == "1h" {
		unitOfTime = "1H"
	} else if payload.Interval == "6h" {
		unitOfTime = "6H"
	} else if payload.Interval == "12h" {
		unitOfTime = "12H"
	} else if payload.Interval == "1d" {
		unitOfTime = "1D"
	} else if payload.Interval == "1w" {
		unitOfTime = "1W"
	}  else if payload.Interval == "1M"{
		unitOfTime = "1MO"
	}

	return TradingChart{trdconst.TRADINGCHART, daos.GetChartList(payload.Topic, unitOfTime)}
}