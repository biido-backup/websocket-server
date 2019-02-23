package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func CreateClient(brokers []string) sarama.Client {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Error("CREATE SARAMA CLIENT ERROR:", err)
	}

	return client
}

func GetOffsetPartition(client sarama.Client, topic string, partition int32) int64 {
	offset, err := client.GetOffset(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("GET OFFSET", topic, " ERROR:", err)
	}

	return offset
}

func CreateConsumerPartition(client sarama.Client, topic string, partition int32, offset int64) sarama.PartitionConsumer {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Error("CREATE CONSUMER", topic, " ERROR:", err)
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Error("CREATE CONSUMER PARTITION", topic, " ERROR:", err)
	}

	return partitionConsumer
}