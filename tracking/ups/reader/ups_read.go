package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
)

var (
	// kafka
	kafkaBrokerUrl     string
	kafkaVerbose       bool
	kafkaTopic         string
	kafkaConsumerGroup string
	kafkaClientId      string
)

type TrackingUpdate struct {
	Index          int    `json:"index"`
	Vendor         string `json:"vendor"`
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
}

func main() {
	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:19092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "shipmentUpdates", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.Parse()

	// make a new consumer that consumes from topic-A
	config := kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerUrl,
		"group.id":          kafkaConsumerGroup,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Error().Msgf("error while receiving message: %s", err.Error())
		return
	}
	defer consumer.Close()

	consumer.Subscribe(kafkaTopic, nil)

	for {
		m, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Error().Msgf("error while receiving message: %s", err.Error())
			continue
		}

		var trackingUpdate TrackingUpdate
		err = json.Unmarshal(m.Value, &trackingUpdate)
		if err != nil {
			log.Error().Msgf("error while unmarshalling: %s", err.Error())
			panic(err)
		}

		fmt.Printf("message at topic/partition/offset %v/%v/%v\n", m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		fmt.Printf("%+v\n", trackingUpdate)
	}
}
