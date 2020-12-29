package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jayxuchen/ups_tracker/kafka"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
)

var logger = log.With().Str("pkg", "main").Logger()

var (
	listenAddrAPI string

	// kafka
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaClientID  string
	kafkaTopic     string
)

//TrackingUpdate describes message from Vendor API to Kafka
type TrackingUpdate struct {
	Vendor         string `json:"vendor"`
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
}

func main() {
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092,localhost:29092,localhost:39092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaClientID, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	flag.StringVar(&kafkaTopic, "kafka-topic", "foo", "Kafka topic to push")

	flag.Parse()

	// connect to kafka
	kafkaProducer, err := kafka.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("unable to configure kafka")
		return
	}
	defer kafkaProducer.Close()

	trackingUpdate := TrackingUpdate{
		Vendor:         "ups",
		TrackingNumber: "1ZABC",
		Status:         "delivered"}

	postDataToKafka(trackingUpdate)
}

func postDataToKafka(msg TrackingUpdate) {
	parent := context.Background()
	defer parent.Done()

	formInBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error while marshalling json: %s", err.Error())
		return
	}
	err = kafka.Push(parent, nil, formInBytes)
	if err != nil {
		fmt.Printf("error while push message into kafka: %s", err.Error())
		return
	}
}
