package logging

import (
	"encoding/json"
	"log"
	"time"
	"os"
)

var tags []string
var labels map[string]string
const ecsVersion = "1.5.0"

func init() {
	tagsFromEnv, exists := os.LookupEnv("HABERDASHER_TAGS")
	if !exists {
		tagsFromEnv = "[]"
	}
	labelsFromEnv, exists := os.LookupEnv("HABERDASHER_LABELS")
	if !exists {
		labelsFromEnv = "{}"
	}
	err := json.Unmarshal([]byte(tagsFromEnv), &tags)
	if err != nil {
		log.Fatal("HABERDASHER_TAGS must be a JSON array of strings")
	}
	err = json.Unmarshal([]byte(labelsFromEnv), &labels)
	if err != nil {
		log.Fatal("HABERDASHER_LABELS must be a JSON object of strings")
	}
}

// An Emitter defines how to ship a log message to a log service.
type Emitter interface {
	HandleLogMessage(jsonBytes []byte) (error)
}

// A Message is a structured log message
type Message struct {
	ECSVersion string `json:"ecs.version"`
	Timestamp time.Time `json:"@timestamp"`
	Labels map[string]string `json:"labels"`
	Tags []string `json:"tags"`
	Message string `json:"message"`
}

// Emitters is the registry of Emitter implementers
var Emitters = make(map[string]Emitter)

// Register will make note of new types of Emitters
func Register(emitterType string, emitter Emitter) {
	Emitters[emitterType] = emitter
}

// Emit is launched as a goroutine for individual log lines to be sent
// concurrently.
func Emit(emitter Emitter, logMessage string) {
	m := Message{ecsVersion, time.Now(), labels, tags, logMessage}
	jsonBytes, _ := json.Marshal(m)
	if err := emitter.HandleLogMessage(jsonBytes); err != nil {
		log.Println("Error emitting message:", jsonBytes, err)
	}
}