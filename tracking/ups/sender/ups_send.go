package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	"strings"
)

var logger = log.With().Str("pkg", "main").Logger()
var writer *kafka.Writer

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
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err = writer.WriteMessages(parent, message)

	if err != nil {
		fmt.Printf("error while push message into kafka: %s", err.Error())
		return
	}
}

//Configure configures kafka writer
func Configure(kafkaBrokerUrls []string, clientID string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}

	config := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}
