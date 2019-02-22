package daos

import (
	"fmt"
	"github.com/linkedin/goavro"
)

var codecTradingChart *goavro.Codec

type TradingChart struct {
	High 		float64
	Low 		float64
	Open 		float64
	Close 		float64
	Volume 		float64
	Mean 		float64
	Time 		int64
}

func init()  {
	recordSchemaJSON := `{
		"type": "record",
		"name": "TradingChart",
		"fields": [
			{"name": "high", "type": "double"},
			{"name": "low", "type": "double"},
			{"name": "open", "type": "double"},
			{"name": "close", "type": "double"},
			{"name": "volume", "type": "double"},
			{"name": "mean", "type": "double"},
			{"name": "time", "type": "long"}
		]
	}`
	codecTradingChart, _ = goavro.NewCodec(recordSchemaJSON)
}

func TradingChartFromJSON(msg []byte) TradingChart {
	native, _, err := codecTradingChart.NativeFromTextual(msg)
	if err != nil {
		fmt.Println(err)
	}

	tradingChartMap := native.(map[string]interface{})

	var tradingChart TradingChart
	tradingChart.High = tradingChartMap["high"].(float64)
	tradingChart.Low = tradingChartMap["low"].(float64)
	tradingChart.Open = tradingChartMap["open"].(float64)
	tradingChart.Close = tradingChartMap["close"].(float64)
	tradingChart.Volume = tradingChartMap["volume"].(float64)
	tradingChart.Mean = tradingChartMap["mean"].(float64)
	tradingChart.Time = tradingChartMap["time"].(int64)

	return tradingChart
}