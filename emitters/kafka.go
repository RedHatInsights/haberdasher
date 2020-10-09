package emitters

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/RedHatInsights/haberdasher/logging"
)

var producer *kafka.Writer
var topic string
type kafkaEmitter struct{}

func init() {
	var emitter kafkaEmitter
	logging.Register("kafka", emitter)
}

// If the Kafka emitter is activated, create a new Producer and spawn a
// goroutine to note any errors.
func (e kafkaEmitter) Setup() {
	bootstrapServers, exists := os.LookupEnv("HABERDASHER_KAFKA_BOOTSTRAP")
	if !exists {
		log.Fatal("To use Haberdasher with Kafka, HABERDASHER_KAFKA_BOOTSTRAP must be set to your bootstrap servers")
	}

	topic, exists = os.LookupEnv("HABERDASHER_KAFKA_TOPIC")
	if !exists {
		log.Fatal("To use Haberdasher with Kafka, HABERDASHER_KAFKA_TOPIC must be set to your logging topic")
	}

	producer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: strings.Split(bootstrapServers, ","),
		Topic: topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// HandleLogMessage ships the log message to Kafka
// TODO: Do we need to re-establish a connection if the Producer connection closes?
func (e kafkaEmitter) HandleLogMessage(jsonSerializeable interface{}) (error) {
	jsonBytes, err := json.Marshal(jsonSerializeable)
	if err != nil {
		err = producer.WriteMessages(
			context.Background(), 
			kafka.Message{
				Value: jsonBytes,
			},
		)
	}
	return err
}

// We don't want any buffered messages to get lost if we shut down, so we wait
// to allow it to exit.
func (e kafkaEmitter) Cleanup() (error) {
	return producer.Close()
}