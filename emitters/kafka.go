package emitters

import (
	"fmt"
	"log"
	"os"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/RedHatInsights/haberdasher/logging"
)

var producer *kafka.Producer
var topic string
type kafkaEmitter struct{}

func init() {
	var emitter kafkaEmitter
	logging.Register("kafka", emitter)
}

// If the Kafka emitter is activated, create a new Producer and spawn a
// goroutine to note any errors.
func (e kafkaEmitter) Setup() {
	var err error
	bootstrapServers, exists := os.LookupEnv("HABERDASHER_KAFKA_BOOTSTRAP")
	if !exists {
		log.Fatal("To use Haberdasher with Kafka, HABERDASHER_KAFKA_BOOTSTRAP must be set to your bootstrap servers")
	}

	topic, exists = os.LookupEnv("HABERDASHER_KAFKA_TOPIC")
	if !exists {
		log.Fatal("To use Haberdasher with Kafka, HABERDASHER_KAFKA_TOPIC must be set to your logging topic")
	}

	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		log.Fatal("Error creating Kafka producer", err)
	}

	go func() {
		for e:= range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Println("Kafka log delivery failure:", ev.TopicPartition)
				}
			}
		}
	}()
}

// HandleLogMessage ships the log message to Kafka
// TODO: Do we need to re-establish a connection if the Producer connection closes?
func (e kafkaEmitter) HandleLogMessage(jsonBytes []byte) (error) {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: jsonBytes,
	}, nil)
}

// We don't want any buffered messages to get lost if we shut down, so we wait
// to allow it to exit. The wait is 9 seconds, expecting that the container
// orchestrator will wait 10 seconds for shutdown.
func (e kafkaEmitter) Cleanup() (error) {
	if messagesRemaining := producer.Flush(9*1000); messagesRemaining > 0 {
		return fmt.Errorf("Failed to flush completely. %d messages still in buffer", messagesRemaining)
	}
	return nil
}