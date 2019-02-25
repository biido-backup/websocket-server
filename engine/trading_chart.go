package engine

import (
	"encoding/json"
	"github.com/Shopify/sarama"
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

	unitOfTimeList := []string{"1MO", "1W", "1D", "12H", "6H", "1H", "30M", "15M", "5M", "1M"}
	consumerRateUnitTimeMap := make(map[string]map[string]sarama.PartitionConsumer)
	for _, tradingRate := range tradingRateList {
		daos.InitRateTradingChart(tradingRate.StringDash())
		consumerRateUnitTimeMap[tradingRate.StringDash()] = make(map[string]sarama.PartitionConsumer)
		for _, unitOfTime := range unitOfTimeList {
			daos.InitUnitOfTimeTradingChart(tradingRate.StringDash(), unitOfTime)
			consumerRateUnitTimeMap[tradingRate.StringDash()][unitOfTime] = loadTradingChart(tradingRate, unitOfTime, quantity)
		}
	}

	for _, tradingRate := range tradingRateList {
		for _, unitOfTime := range unitOfTimeList {
			go maintainTradingChart(tradingRate, unitOfTime, quantity, consumerRateUnitTimeMap[tradingRate.StringDash()][unitOfTime])
		}
	}
}

func loadTradingChart(rate daos.Rate, unitOfTime string, quantity int64) sarama.PartitionConsumer {
	var consumer sarama.PartitionConsumer

	kafkaTopic := viper.GetString("kafka.prefix.chart") + unitOfTime + "." + rate.StringDash()
	brokers := viper.GetStringSlice("kafka.websocket.brokers")
	client := kafka.CreateClient(brokers)

	chartList := make([]daos.Chart, 0, quantity)
	offset := kafka.GetOffsetPartition(client, kafkaTopic, 0)

	if offset > 0 {
		offsetToConsume := offset - quantity
		if offsetToConsume < 0 {
			offsetToConsume = 0
		}

		consumer = kafka.CreateConsumerPartition(client, kafkaTopic, 0, offsetToConsume)
		for {
			msg := <- consumer.Messages()
			chart := daos.ChartFromJSON(msg.Value)
			chartList = append(chartList, chart)

			if unitOfTime == "1M" {
				daos.CalculateNextChart(rate.StringDash(), chart)
			}

			if msg.Offset == offset - 1 {
				if unitOfTime != "1M" {
					daos.SetNextChart(rate.StringDash(), unitOfTime, chart.Time)
				}
				break
			}
		}

		daos.SetChartList(rate.StringDash(), unitOfTime, chartList)
		//fmt.Println(daos.TRDChart)
	} else {
		consumer = kafka.CreateConsumerPartition(client, kafkaTopic, 0, 0)
	}

	return consumer
}

func maintainTradingChart(rate daos.Rate, unitOfTime string, quantity int64, consumer sarama.PartitionConsumer) {
	for {
		msg := <- consumer.Messages()
		chart := daos.ChartFromJSON(msg.Value)
		daos.InsertChart(rate.StringDash(), unitOfTime, chart, quantity)
		if unitOfTime == "1M" {
			daos.CalculateNextChart(rate.StringDash(), chart)
		}
		//fmt.Println(daos.TRDChart)

		tradingChart := trading.TradingChart{trdconst.TRADINGCHART, []daos.Chart{daos.ChartFromJSON(msg.Value)}}
		tradingChartJson, _ := json.Marshal(tradingChart)

		if unitOfTime == "1M" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "1m", string(tradingChartJson))
		} else if unitOfTime == "5M" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "5m", string(tradingChartJson))
		} else if unitOfTime == "15M" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "15m", string(tradingChartJson))
		} else if unitOfTime == "30M" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "30m", string(tradingChartJson))
		} else if unitOfTime == "1H" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "1h", string(tradingChartJson))
		} else if unitOfTime == "6H" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "6h", string(tradingChartJson))
		} else if unitOfTime == "12H" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "12h", string(tradingChartJson))
		} else if unitOfTime == "1D" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "1d", string(tradingChartJson))
		} else if unitOfTime == "1W" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "1w", string(tradingChartJson))
		}  else if unitOfTime == "1MO" {
			websocket.BroadcastMessageWithInterval(rate.StringDash(), "1M", string(tradingChartJson))
		}
	}
}