package daos

import (
	"fmt"
	"github.com/linkedin/goavro"
)

var codecChart *goavro.Codec

type Chart struct {
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
		"name": "Chart",
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
	codecChart, _ = goavro.NewCodec(recordSchemaJSON)
}

func ChartFromJSON(msg []byte) (chart Chart) {
	native, _, err := codecChart.NativeFromTextual(msg)
	if err != nil {
		fmt.Println(err)
	}
	chartMap := native.(map[string]interface{})

	chart.High = chartMap["high"].(float64)
	chart.Low = chartMap["low"].(float64)
	chart.Open = chartMap["open"].(float64)
	chart.Close = chartMap["close"].(float64)
	chart.Volume = chartMap["volume"].(float64)
	chart.Mean = chartMap["mean"].(float64)
	chart.Time = chartMap["time"].(int64)

	return chart
}