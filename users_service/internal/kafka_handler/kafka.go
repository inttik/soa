package kafkahandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter() (*KafkaWriter, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "registration",
	})
	if writer == nil {
		return nil, errors.New("can't connect to kafka")
	}
	return &KafkaWriter{
		writer: writer,
	}, nil
}

func (w *KafkaWriter) Close() {
	if w == nil {
		return
	}
	w.writer.Close()
}

func (w *KafkaWriter) NewClient(userId uuid.UUID) error {
	if w == nil {
		return nil
	}
	registrationTime := time.Now().Format("2006-01-02 15:04:05")

	message := fmt.Sprintf("user_id: %s, reg_time: %s", userId.String(), registrationTime)

	return w.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(userId.String()),
			Value: []byte(message),
		})
}
