package engine

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/kafka"
	"websocket-server/util/websocket"
)

func ProcessTradingChart(tradingRateList []daos.Rate) {
	daos.CreateTradingChart()
	quantity := viper.GetInt64("tradingchart.quantity")

	unitOfTimeList := []string{"MINUTE", "HOUR", "DAY", "MONTH"}
	for _, tradingRate := range tradingRateList {
		for _, unitOfTime := range unitOfTimeList {
			maintainTradingChart(tradingRate, unitOfTime, quantity)
		}
	}
}

func maintainTradingChart(rate daos.Rate, unitOfTime string, quantity int64) {
	kafkaTopic := viper.GetString("kafka.prefix.chart") + rate.StringDash() + "-" + unitOfTime
	brokers := viper.GetStringSlice("kafka.websocket.brokers")
	client := kafka.CreateClient(brokers)

	chartList := make([]daos.Chart, 0, quantity)
	offset := kafka.GetOffsetPartition(client, kafkaTopic, 0)
	if offset > 0 {
		offsetToConsume := offset - quantity
		if offsetToConsume < 0 {
			offsetToConsume = 0
		}

		partitionConsumer := kafka.CreateConsumerPartition(client, kafkaTopic, 0, offsetToConsume)
		for {
			msg := <- partitionConsumer.Messages()
			chart := daos.ChartFromJSON(msg.Value)

			chartList = append(chartList, chart)
			if msg.Offset == offset - 1 {
				break
			}
		}

		daos.SetChartList(rate.StringDash(), unitOfTime, chartList)

		for {
			msg := <- partitionConsumer.Messages()
			daos.InsertChart(rate.StringDash(), unitOfTime, daos.ChartFromJSON(msg.Value), quantity)

			tradingChart := trading.TradingChart{trdconst.TRADINGCHART, []daos.Chart{daos.ChartFromJSON(msg.Value)}}
			tradingChartJson, _ := json.Marshal(tradingChart)
			fmt.Println(string(tradingChartJson))

			if unitOfTime == "MINUTE" {
				websocket.BroadcastMessageWithInterval(rate.StringDash(), "1m", string(tradingChartJson))
			} else if unitOfTime == "HOUR" {
				websocket.BroadcastMessageWithInterval(rate.StringDash(), "1h", string(tradingChartJson))
			} else if unitOfTime == "DAY" {
				websocket.BroadcastMessageWithInterval(rate.StringDash(), "1d", string(tradingChartJson))
				websocket.BroadcastMessageWithInterval(rate.StringDash(), "1w", string(tradingChartJson))
			} else if unitOfTime == "MONTH" {
				websocket.BroadcastMessageWithInterval(rate.StringDash(), "1M", string(tradingChartJson))
			}
		}
	}
}