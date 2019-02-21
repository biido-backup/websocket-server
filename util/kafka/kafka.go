package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"log"
	"sync"
	"websocket-server/const/trd"
	"websocket-server/daos"
	"websocket-server/daos/trading"
	"websocket-server/util/websocket"
)

func GetCandleStick(rate, unitOfTime string) (tradingChartJSONList []string, offset int64) {
	tradingChartJSONList = make([]string, 0, 5)
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
		offsetToConsume := offset - 5
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
			tradingChartJSONList = append(tradingChartJSONList, string(msg.Value))
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
		size := len(tradingCharts.ListMap[rate][unitOfTime])
		if size >= 5 {
			tradingCharts.ListMap[rate][unitOfTime] = tradingCharts.ListMap[rate][unitOfTime][1:len(tradingCharts.ListMap[rate][unitOfTime])]
		}
		tradingCharts.ListMap[rate][unitOfTime] = append(tradingCharts.ListMap[rate][unitOfTime], string(msg.Value))
		tradingCharts.OffsetMap[rate][unitOfTime] = msg.Offset

		tradingChart := trading.TradingChart{trdconst.TRADINGCHART, []string{string(msg.Value)}}
		tradingChartJson, _ := json.Marshal(tradingChart)
		//fmt.Println(string(tradingChartJson))
		websocket.BroadcastMessage(rate, string(tradingChartJson))
		mutex.Unlock()
	}
}