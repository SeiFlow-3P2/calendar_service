package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	Produce(ctx context.Context, msg interface{}, topic string, key string, timeout time.Duration) error
	Close() error
}

type saramaProducer struct {
	producer sarama.SyncProducer
}

func NewProducer(address []string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama sync producer: %w", err)
	}
	log.Printf("Kafka producer connected to brokers: %v", address)
	return &saramaProducer{producer: producer}, nil
}

func (p *saramaProducer) Produce(ctx context.Context, msg interface{}, topic string, key string, timeout time.Duration) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(jsonMsg),
	}

	partition, offset, err := p.producer.SendMessage(producerMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka topic %s: %w", topic, err)
	}

	log.Printf("Message sent to Kafka topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

func (p *saramaProducer) Close() error {
	if p.producer != nil {
		log.Println("Closing Kafka producer...")
		return p.producer.Close()
	}
	return nil
}
