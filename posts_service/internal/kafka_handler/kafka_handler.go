package kafkahandler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaHandler struct {
	likes    *kafka.Writer
	watches  *kafka.Writer
	comments *kafka.Writer
}

func NewKafkaHandler() (*KafkaHandler, error) {
	likes := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "likes",
	})
	if likes == nil {
		return nil, errors.New("can't connect to kafka: likes")
	}

	watches := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "watches",
	})
	if watches == nil {
		return nil, errors.New("can't connect to kafka: watches")
	}

	comments := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "comments",
	})
	if comments == nil {
		return nil, errors.New("can't connect to kafka: watches")
	}

	return &KafkaHandler{
		likes:    likes,
		watches:  watches,
		comments: comments,
	}, nil
}

func (w *KafkaHandler) Close() {
	if w == nil {
		return
	}

	w.likes.Close()
	w.watches.Close()
	w.comments.Close()
}

func (w *KafkaHandler) NewLike(userId uuid.UUID, post_id uuid.UUID) error {
	if w == nil {
		return nil
	}

	likeTime := time.Now().Format("2006-01-02 15:04:05")

	message := fmt.Sprintf("user_id: %s, post_id: %s, like_time: %s", userId.String(), post_id.String(), likeTime)

	return w.likes.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
}

func (w *KafkaHandler) NewWatch(userId uuid.UUID, post_id uuid.UUID) error {
	if w == nil {
		return nil
	}

	watchTime := time.Now().Format("2006-01-02 15:04:05")

	message := fmt.Sprintf("user_id: %s, post_id: %s, watch_time: %s", userId.String(), post_id.String(), watchTime)

	return w.watches.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
}

func (w *KafkaHandler) NewComment(userId uuid.UUID, post_id uuid.UUID, text string) error {
	if w == nil {
		return nil
	}

	commentTime := time.Now().Format("2006-01-02 15:04:05")

	message := fmt.Sprintf("user_id: %s, post_id: %s, comment_time: %s\ntext: %s", userId.String(), post_id.String(), commentTime, text)

	return w.comments.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
}
