package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"log"
	"sync"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/websocket"
)

func GetTradingChart(rate, unitOfTime string) (tradingChartJSONList []daos.TradingChart, offset int64) {
	tradingChartJSONList = make([]daos.TradingChart, 0, 150)
	kafkaTopic := rate + ".TRADING_CHART." + unitOfTime

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	webappBrokers := viper.GetStringSlice("kafka.webapp.brokers")

	client, err := sarama.NewClient(webappBrokers, config)
	if err != nil {
		log.Fatal("Unable to connect client to kafka trading chart")
	}

	offset, err = client.GetOffset(kafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed get offset order book")
	}

	if offset > 0 {
		offsetToConsume := offset - 150
		if offsetToConsume < 0 {
			offsetToConsume = 0
		}

		consumer, err := sarama.NewConsumerFromClient(client)
		if err != nil {
			log.Fatal("Unable to connect consumer to kafka trading chart")
		}

		partitionConsumer, err := consumer.ConsumePartition(kafkaTopic, 0, offsetToConsume)
		if err != nil {
			log.Fatal("Failed to consume partition kafka trading chart")
		}

		for {
			msg := <- partitionConsumer.Messages()
			tradingChart := daos.TradingChartFromJSON(msg.Value)

			length := len(tradingChartJSONList)
			if length > 0 {
				lowerBound := 0
				upperBound := length - 1
				curIn := 0

				for {
					curIn = (upperBound + lowerBound) / 2
					if tradingChartJSONList[curIn].Time < tradingChart.Time {
						lowerBound = curIn + 1
						if lowerBound > upperBound {
							curIn = curIn + 1
							break
						}
					} else if tradingChartJSONList[curIn].Time > tradingChart.Time {
						upperBound = curIn - 1
						if lowerBound > upperBound {
							break
						}
					} else {
						break
					}
				}

				if curIn >= length {
					tradingChartJSONList = append(tradingChartJSONList, tradingChart)
				} else {
					tradingChartJSONList = append(tradingChartJSONList[:curIn+1], tradingChartJSONList[curIn:]...)
					tradingChartJSONList[curIn] = tradingChart
				}
			} else {
				tradingChartJSONList = append(tradingChartJSONList, tradingChart)
			}

			if msg.Offset == offset - 1 {
				break
			}
		}
		partitionConsumer.Close()
		consumer.Close()
	}
	client.Close()

	return tradingChartJSONList, offset
}

func MaintenanceTradingChartArray(rate string, unitOfTime string, tradingCharts *daos.TradingCharts, mutex *sync.Mutex) {
	kafkaTopic := rate + ".TRADING_CHART." + unitOfTime

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	webappBrokers := viper.GetStringSlice("kafka.webapp.brokers")

	client, err := sarama.NewClient(webappBrokers, config)
	if err != nil {
		log.Fatal("Unable to connect client to kafka trading chart")
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal("Unable to connect consumer to kafka trading chart")
	}


	partitionConsumer, err := consumer.ConsumePartition(kafkaTopic, 0, tradingCharts.OffsetMap[rate][unitOfTime])
	if err != nil {
		log.Fatal("Failed to consume partition kafka trading chart")
	}

	for {
		msg := <- partitionConsumer.Messages()

		mutex.Lock()
		if len(tradingCharts.ListMap[rate][unitOfTime]) >= 150 {
			tradingCharts.ListMap[rate][unitOfTime] = tradingCharts.ListMap[rate][unitOfTime][1:len(tradingCharts.ListMap[rate][unitOfTime])]
		}
		tradingCharts.ListMap[rate][unitOfTime] = append(tradingCharts.ListMap[rate][unitOfTime], daos.TradingChartFromJSON(msg.Value))
		tradingCharts.OffsetMap[rate][unitOfTime] = msg.Offset


		tradingChart := trading.TradingChart{trdconst.TRADINGCHART, []daos.TradingChart{daos.TradingChartFromJSON(msg.Value)}}
		tradingChartJson, _ := json.Marshal(tradingChart)
		fmt.Println(string(tradingChartJson))
		websocket.BroadcastMessage(rate, string(tradingChartJson))
		mutex.Unlock()
	}
}