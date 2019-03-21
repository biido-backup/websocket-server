package engine

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/kafka"
	"websocket-server/util/logger"
	"websocket-server/util/websocket"
)

var log = logger.CreateLog("engine")

func ProcessTradingChart(tradingRateList []daos.Rate) {
	daos.CreateTradingChart()
	size := viper.GetInt64("tradingchart.size")

	//tradingRateList = []daos.Rate{{"BION", "ETH"}}
	//unitOfTimeList := []string{"1M"}
	unitOfTimeList := []string{"1M", "5M", "15M", "30M", "1H", "6H", "12H", "1D", "1W", "1MO"}
	consumerRateUnitTimeMap := make(map[string]map[string]sarama.PartitionConsumer)
	for _, tradingRate := range tradingRateList {
		daos.InitRateTradingChart(tradingRate.StringDash())
		consumerRateUnitTimeMap[tradingRate.StringDash()] = make(map[string]sarama.PartitionConsumer)
		for _, unitOfTime := range unitOfTimeList {
			daos.InitUnitOfTimeTradingChart(tradingRate.StringDash(), unitOfTime)
			consumerRateUnitTimeMap[tradingRate.StringDash()][unitOfTime] = loadTradingChart(tradingRate, unitOfTime, size)
		}
	}

	for _, tradingRate := range tradingRateList {
		for _, unitOfTime := range unitOfTimeList {
			go maintainTradingChart(tradingRate, unitOfTime, size, consumerRateUnitTimeMap[tradingRate.StringDash()][unitOfTime])
		}
	}
}

func loadTradingChart(rate daos.Rate, unitOfTime string, size int64) sarama.PartitionConsumer {
	var consumer sarama.PartitionConsumer

	kafkaTopic := viper.GetString("kafka.prefix.chart") + unitOfTime + "." + rate.StringDash()
	brokers := viper.GetStringSlice("kafka.websocket.brokers")
	client := kafka.CreateClient(brokers)

	chartList := make([]daos.Chart, 0, size)
	offset := kafka.GetOffsetPartition(client, kafkaTopic, 0)

	if offset > 0 {
		offsetToConsume := offset - size
		if offsetToConsume < 0 {
			offsetToConsume = 0
		}

		consumer = kafka.CreateConsumerPartition(client, kafkaTopic, 0, offsetToConsume)
		for {
			msg := <- consumer.Messages()
			chart := daos.ChartFromJSON(msg.Value)

			chartList = append(chartList, chart)
			if msg.Offset == offset - 1 {
				break
			}
		}

		daos.SetChartList(rate.StringDash(), unitOfTime, chartList)
		//log.Debug(kafkaTopic, daos.GetChartList(rate.StringDash(), unitOfTime))
	} else {
		consumer = kafka.CreateConsumerPartition(client, kafkaTopic, 0, 0)
	}

	return consumer
}

func maintainTradingChart(rate daos.Rate, unitOfTime string, size int64, consumer sarama.PartitionConsumer) {
	for {
		msg := <- consumer.Messages()
		daos.InsertChart(rate.StringDash(), unitOfTime, daos.ChartFromJSON(msg.Value), size)

		tradingChart := trading.TradingChart{trdconst.TRADINGCHART, []daos.Chart{daos.ChartFromJSON(msg.Value)}}
		tradingChartJson, _ := json.Marshal(tradingChart)
		//log.Debug(string(tradingChartJson))
		//log.Debug(rate.StringDash(), unitOfTime, daos.GetChartList(rate.StringDash(), unitOfTime))

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