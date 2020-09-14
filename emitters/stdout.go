package emitters

import (
	"fmt"
	"github.com/RedHatInsights/haberdasher/logging"
)

type stdoutEmitter struct{}

func init() {
	var emitter stdoutEmitter
	logging.Register("stdout", emitter)
}

// HandleLogMessage prints the json body to stdout
func (e stdoutEmitter) HandleLogMessage(jsonBytes []byte) (error) {
	fmt.Println(string(jsonBytes))
	return nil
}

func (e stdoutEmitter) Cleanup() (error) {
	return nil
}
