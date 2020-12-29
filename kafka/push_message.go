package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"time"
)

//Push pushes message to topic
func Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(parent, message)
}
