package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
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
	Vendor         string `json:"vendor"`
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
}

func main() {
	flag.StringVar(&kafkaBrokerUrl, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic. Only one topic per worker.")
	flag.StringVar(&kafkaConsumerGroup, "kafka-consumer-group", "consumer-group", "Kafka consumer group")
	flag.StringVar(&kafkaClientId, "kafka-client-id", "my-client-id", "Kafka client id")

	flag.Parse()

	brokers := strings.Split(kafkaBrokerUrl, ",")

	// make a new reader that consumes from topic-A
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaClientId,
		Topic:           kafkaTopic,
		Partition:       0,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
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

		fmt.Printf("message at topic/partition/offset %v/%v/%v\n", m.Topic, m.Partition, m.Offset)
		fmt.Printf("%+v\n", trackingUpdate)
	}
}
