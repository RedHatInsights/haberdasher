package logging

import (
	"encoding/json"
	"log"
)

// An Emitter defines how to ship a log message to a log service.
type Emitter interface {
	HandleLogMessage(jsonBytes []byte) (error)
}

// A Message is a structured log message
type Message struct {
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
	m := Message{logMessage}
	jsonBytes, _ := json.Marshal(m)
	if err := emitter.HandleLogMessage(jsonBytes); err != nil {
		log.Println("Error emitting message:", jsonBytes, err)
	}
}