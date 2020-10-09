package emitters

import (
	"fmt"
	"encoding/json"
	"os"
	"github.com/RedHatInsights/haberdasher/logging"
)

type stderrEmitter struct{}

func init() {
	var emitter stderrEmitter
	logging.Register("stderr", emitter)
}

func (e stderrEmitter) Setup() {}

func (e stderrEmitter) HandleLogMessage(jsonSerializeable interface{}) (error) {
	var jsonBytes []byte
	var err error
	prettyPrint := os.Getenv("HABERDASHER_STDERR_PRETTY")
	if prettyPrint != "" {
		jsonBytes, err = json.MarshalIndent(jsonSerializeable, "", "    ")
	} else {
		jsonBytes, err = json.Marshal(jsonSerializeable)
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, string(jsonBytes))
	return nil
}

func (e stderrEmitter) Cleanup() (error) {
	return nil
}
