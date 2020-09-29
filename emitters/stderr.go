package emitters

import (
	"fmt"
	"os"
	"github.com/RedHatInsights/haberdasher/logging"
)

type stderrEmitter struct{}

func init() {
	var emitter stderrEmitter
	logging.Register("stderr", emitter)
}

func (e stderrEmitter) Setup() {}

func (e stderrEmitter) HandleLogMessage(jsonBytes []byte) (error) {
	fmt.Fprintln(os.Stderr, string(jsonBytes))
	return nil
}

func (e stderrEmitter) Cleanup() (error) {
	return nil
}
