package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/namsral/flag"
	"github.com/rs/zerolog/log"
	//	"time"
)

var logger = log.With().Str("pkg", "main").Logger()
var producer *kafka.Producer

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
	Index          int    `json:"index"`
	Vendor         string `json:"vendor"`
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
}

func main() {
	flag.StringVar(&kafkaBrokerURL, "kafka-brokers", "localhost:19092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&kafkaClientID, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	flag.StringVar(&kafkaTopic, "kafka-topic", "shipmentUpdates", "Kafka topic to push")

	flag.Parse()

	// connect to kafka
	producer, err := Configure(kafkaBrokerURL, kafkaClientID)
	if err != nil {
		logger.Error().Str("error", err.Error()).Msg("unable to configure kafka")
		return
	}
	defer producer.Close()

	trackingUpdate := TrackingUpdate{
		Index:          0,
		Vendor:         "ups",
		TrackingNumber: "1ZABC",
		Status:         "delivered"}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for i := 0; i < 10; i++ {
		trackingUpdate.Index = i
		postDataToKafka(trackingUpdate, kafkaTopic)
		//		time.Sleep(time.Second)
	}
	producer.Flush(15 * 1000)
}

func postDataToKafka(msg TrackingUpdate, topic string) {

	formInBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error while marshalling json: %s", err.Error())
		return
	}
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          formInBytes,
	}

	err = producer.Produce(&message, nil)

	if err != nil {
		fmt.Printf("error while push message into kafka: %s", err.Error())
		return
	}
}

//Configure configures kafka writer
func Configure(kafkaBrokerUrls string, clientID string) (w *kafka.Producer, err error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBrokerUrls,
		"client.id":         clientID,
		"acks":              "all"})
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
	}
	producer = p
	return p, err
}
