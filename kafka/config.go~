package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

var writer *kafka.Writer

//Configure configures kafka connection for kafka writer
func Configure(kafkaBrokerUrls []string, clientID string, topic string) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}

	config := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        topic,
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}
