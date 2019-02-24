package kafka

import (
	"github.com/Shopify/sarama"
)



func CreateProducer(brokers []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = false
	config.Producer.MaxMessageBytes = 50000000

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Error("Error create producer", err)
	}

	return producer
}

func Publish(message string, producer sarama.AsyncProducer, topic string) {
	msg := &sarama.ProducerMessage {
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	producer.Input() <- msg
}