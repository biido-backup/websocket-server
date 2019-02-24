package kafka

import (
	"github.com/Shopify/sarama"
	"websocket-server/util/logger"
)

var log = logger.CreateLog("kafka")

func CreateClient(brokers []string) sarama.Client {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Error("Error create client", err)
	}

	return client
}

func GetOffsetPartition(client sarama.Client, topic string, partition int32) int64 {
	offset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Error("Error get offset for topic : "+topic, err)
	}

	return offset
}

func CreateConsumerPartition(client sarama.Client, topic string, partition int32, offset int64) sarama.PartitionConsumer {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Error("Error create consmer for topic : "+topic, err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Error("Error create consmer partitions for topic : "+topic, err)
	}

	return partitionConsumer
}