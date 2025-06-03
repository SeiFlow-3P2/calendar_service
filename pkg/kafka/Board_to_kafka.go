package kafka

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	kafkaBootstrapServers = "kafka-broker:9092"
	groupID               = "board-event-processor"
	topic                 = "board-events"
)

type BoardEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  kafkaBootstrapServers,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		fmt.Printf("Ошибка создания потребителя: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()

	// Подписка на топик
	if err := consumer.Subscribe(topic, nil); err != nil {
		fmt.Printf("Ошибка подписки на топик: %v\n", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigchan:
			fmt.Printf("Получен сигнал %v: завершение работы\n", sig)
			return
		default:
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				fmt.Printf("Ошибка потребления: %v (%v)\n", err, msg)
				continue
			}

			if err := processMessage(msg); err != nil {
				fmt.Printf("Ошибка обработки сообщения: %v\n", err)
				continue
			}

			// Подтверждение обработки сообщения
			if _, err := consumer.CommitMessage(msg); err != nil {
				fmt.Printf("Ошибка подтверждения сообщения: %v\n", err)
			}
		}
	}
}

func processMessage(msg *kafka.Message) error {
	var event BoardEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	switch event.Type {
	case "task-created":
		return handleTaskCreated(event.Payload)
	case "task-updated":
		return handleTaskUpdated(event.Payload)
	case "task-deleted":
		return handleTaskDeleted(event.Payload)
	default:
		return fmt.Errorf("неизвестный тип события: %s", event.Type)
	}
}

func handleTaskCreated(payload json.RawMessage) error {

	fmt.Println("Обработка создания задачи...")
	return nil
}

func handleTaskUpdated(payload json.RawMessage) error {

	fmt.Println("Обработка обновления задачи...")
	return nil
}

func handleTaskDeleted(payload json.RawMessage) error {

	fmt.Println("Обработка удаления задачи...")
	return nil
}
